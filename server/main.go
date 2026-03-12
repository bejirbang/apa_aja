package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"apa_aja/database"
	pb "apa_aja/proto"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

type ctxKey string

const ctxUserID ctxKey = "userID"

type googleProfile struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}

func (s *server) LoginWithGoogle(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.AuthResponse, error) {
	if strings.TrimSpace(req.IdToken) == "" {
		return nil, status.Error(codes.InvalidArgument, "id_token wajib diisi")
	}

	profile, err := verifyGoogleIDToken(ctx, req.IdToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token google tidak valid")
	}

	user, err := upsertGoogleUser(profile)
	if err != nil {
		return nil, status.Error(codes.Internal, "gagal menyimpan user")
	}

	token, expiresAt, err := createSession(user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "gagal membuat session")
	}

	log.Println("login google sukses user_id:", user.Id, "expires:", expiresAt.Format(time.RFC3339))

	return &pb.AuthResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	db := database.DB

	var user pb.UserResponse
	err := db.QueryRow("SELECT id, name, age, email, avatar_url FROM user WHERE id = ?", req.Id).
		Scan(&user.Id, &user.Name, &user.Age, &user.Email, &user.AvatarUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println("Request user ID:", req.Id, "-> ditemukan")

	return &user, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	db := database.DB

	res, err := db.Exec("INSERT INTO user (name, age) VALUES (?, ?)", req.Name, req.Age)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:   int32(lastID),
		Name: req.Name,
		Age:  req.Age,
	}, nil
}

func main() {
	database.InitDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authUnaryInterceptor))
	pb.RegisterUserServiceServer(grpcServer, &server{})

	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	go func() {
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(r) || wrappedGrpc.IsGrpcWebSocketRequest(r) || wrappedGrpc.IsAcceptableGrpcCorsRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
				return
			}

			http.NotFound(w, r)
		})

		fmt.Println("gRPC-Web server running on port 8080")
		if err := http.ListenAndServe(":8080", httpHandler); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("gRPC server running on port 50051")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func authUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == "/user.UserService/LoginWithGoogle" {
		return handler(ctx, req)
	}

	token, err := extractBearerToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak ditemukan")
	}

	userID, err := validateSession(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}

	ctx = context.WithValue(ctx, ctxUserID, userID)
	return handler(ctx, req)
}

func extractBearerToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata kosong")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", errors.New("authorization tidak ada")
	}

	value := strings.TrimSpace(values[0])
	if strings.HasPrefix(strings.ToLower(value), "bearer ") {
		return strings.TrimSpace(value[7:]), nil
	}

	return value, nil
}

func validateSession(token string) (int32, error) {
	db := database.DB

	var userID int32
	var expiresAt time.Time
	err := db.QueryRow("SELECT user_id, expires_at FROM session WHERE token = ?", token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		_, _ = db.Exec("DELETE FROM session WHERE token = ?", token)
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func verifyGoogleIDToken(ctx context.Context, idToken string) (*googleProfile, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	if clientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID belum diatur")
	}

	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, err
	}

	profile := &googleProfile{
		Sub: payload.Subject,
	}

	if email, ok := payload.Claims["email"].(string); ok {
		profile.Email = email
	}
	if name, ok := payload.Claims["name"].(string); ok {
		profile.Name = name
	}
	if picture, ok := payload.Claims["picture"].(string); ok {
		profile.Picture = picture
	}
	if profile.Name == "" {
		profile.Name = profile.Email
	}

	if profile.Sub == "" {
		return nil, errors.New("subject kosong")
	}

	return profile, nil
}

func upsertGoogleUser(profile *googleProfile) (*pb.UserResponse, error) {
	db := database.DB

	var user pb.UserResponse
	err := db.QueryRow("SELECT id, name, age, email, avatar_url FROM user WHERE google_sub = ?", profile.Sub).
		Scan(&user.Id, &user.Name, &user.Age, &user.Email, &user.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res, err := db.Exec(
				"INSERT INTO user (name, age, email, google_sub, avatar_url) VALUES (?, ?, ?, ?, ?)",
				profile.Name, 0, profile.Email, profile.Sub, profile.Picture,
			)
			if err != nil {
				return nil, err
			}

			lastID, err := res.LastInsertId()
			if err != nil {
				return nil, err
			}

			user = pb.UserResponse{
				Id:        int32(lastID),
				Name:      profile.Name,
				Age:       0,
				Email:     profile.Email,
				AvatarUrl: profile.Picture,
			}
			return &user, nil
		}
		return nil, err
	}

	user.Name = profile.Name
	user.Email = profile.Email
	user.AvatarUrl = profile.Picture

	_, err = db.Exec("UPDATE user SET name = ?, email = ?, avatar_url = ? WHERE id = ?", user.Name, user.Email, user.AvatarUrl, user.Id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func createSession(userID int32) (string, time.Time, error) {
	token, err := generateToken(32)
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(getSessionTTL())
	_, err = database.DB.Exec("INSERT INTO session (token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func generateToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func getSessionTTL() time.Duration {
	if raw := strings.TrimSpace(os.Getenv("SESSION_TTL_HOURS")); raw != "" {
		if hours, err := strconv.Atoi(raw); err == nil && hours > 0 {
			return time.Duration(hours) * time.Hour
		}
	}
	return 24 * time.Hour
}

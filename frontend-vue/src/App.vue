<template>
  <div class="app-container">
    <div class="card">
      <h1 class="title">ðŸš€ Microservice Demo</h1>
      <p class="subtitle">Contoh implementasi frontend untuk microservice</p>

      <div class="user-section">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Background utama */
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #4facfe, #00f2fe);
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

/* Card utama */
.card {
  background: white;
  padding: 40px;
  border-radius: 16px;
  width: 600px;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

/* Judul */
.title {
  font-size: 32px;
  margin-bottom: 10px;
  color: #333;
}

/* Subjudul */
.subtitle {
  color: #777;
  margin-bottom: 30px;
}

/* Section user */
.user-section {
  background: #f5f7ff;
  padding: 20px;
  border-radius: 10px;
}
</style>

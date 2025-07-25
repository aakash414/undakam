<template>
  <div class="container">
    <div class="test-card">
      <h1>Undakam Recipe Test</h1>
      
      <div class="input-group">
        <label for="subdomain">Recipe Subdomain:</label>
        <input 
          id="subdomain"
          v-model="testSubdomain" 
          placeholder="chicken-curry-nadan-for10-eluppam"
          @keyup.enter="testRecipe"
          :disabled="loading"
        />
        <button 
          @click="testRecipe"
          :disabled="loading || !testSubdomain.trim()"
          class="test-btn"
        >
          {{ loading ? 'Loading...' : 'Test Recipe' }}
        </button>
      </div>

      <div v-if="error" class="error">
        Error: {{ error }}
      </div>

      <div v-if="result" class="result">
        <h3>Recipe</h3>
        <pre>{{ JSON.stringify(result, null, 2) }}</pre>
      </div>

      <div v-if="loading" class="loading">
        Fetching data
      </div>
    </div>
  </div>
</template>

<script setup>
const { fetchRecipe } = useRecipe()
const testSubdomain = ref('chicken-curry-nadan-for10-eluppam')
const result = ref(null)
const loading = ref(false)
const error = ref(null)

const testRecipe = async () => {
  if (!testSubdomain.value.trim()) return
  
  loading.value = true
  error.value = null
  result.value = null
  
  try {
    console
    const data = await fetchRecipe(testSubdomain.value.trim())
    result.value = data
    
    if (!data) {
      error.value = 'no data'
    }
  } catch (err) {
    error.value = err.message || 'Failed to fetch recipe'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  testRecipe()
})
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 1rem;
  font-family: system-ui, -apple-system, sans-serif;
}

.test-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

h1 {
  text-align: center;
  color: #1f2937;
  margin-bottom: 2rem;
  font-size: 2rem;
}

.input-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #374151;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #d1d5db;
  border-radius: 8px;
  font-size: 1rem;
  margin-bottom: 1rem;
  transition: border-color 0.2s;
}

input:focus {
  outline: none;
  border-color: #3b82f6;
}

input:disabled {
  background-color: #f3f4f6;
  cursor: not-allowed;
}

.test-btn {
  width: 100%;
  padding: 0.75rem 1.5rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.test-btn:hover:not(:disabled) {
  background-color: #2563eb;
}

.test-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

.error {
  background-color: #fef2f2;
  color: #dc2626;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #fca5a5;
  margin-bottom: 1rem;
}

.result {
  background-color: #f0fdf4;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #bbf7d0;
}

.result h3 {
  color: #15803d;
  margin-top: 0;
  margin-bottom: 1rem;
}

.result pre {
  background-color: #1f2937;
  color: #f9fafb;
  padding: 1rem;
  border-radius: 6px;
  overflow-x: auto;
  font-size: 0.875rem;
  line-height: 1.5;
}

.loading {
  text-align: center;
  color: #6b7280;
  font-weight: 500;
  padding: 2rem;
}

@media (max-width: 640px) {
  .container {
    margin: 1rem;
    padding: 0.5rem;
  }
  
  .test-card {
    padding: 1.5rem;
  }
  
  h1 {
    font-size: 1.5rem;
  }
}
</style>

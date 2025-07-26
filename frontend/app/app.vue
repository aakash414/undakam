<template>
  <div class="container">
  <div v-if="!result && !loading" class="hero-section">
    <div class="title-container">
      <h1 class="title">Undakam Recipe Maker</h1>
      <p>generate recipe from simple discription</p>
    </div>
    <div class="search-section">
      <div class="search-container">
        <input
          type="text"
          v-model="searchInput"
          placeholder="chicken-curry-nadan-for10"
          class="search-input"
        />
        <button class="search-btn" @click="searchRecipe">Generate Recipe</button>
      </div>
    </div>
  </div>

  <div v-if="loading">
    <div class="loading-section">
      <h2>Loading...</h2>
      <p>Please wait while we fetch your recipe.</p>
    </div>
  </div>

  <div v-if="result && !loading" class="recipe-section">
    <div class="recipe-title">
      <h2>{{ result.name }}</h2>
      <div class="recipe-details">
        <div>
          <span class="recipe-time">Time: {{ result.cookTime }} minutes</span>
          <span class="recipe-servings">Servings: {{ result.servings }}</span>
        </div>
        <div class="recipe-ingredients">
          <h2>Ingredients</h2>
          <ul>
            <li v-for="ingredient in result.ingredients">
               {{ ingredient }} 
            </li>
          </ul>
        </div>
        <div class="recipe-instructions">
          <h2>Instructions</h2>
          <ol>
            <li v-for="step in result.steps" >
              {{ step }}
            </li>
          </ol>
        </div>
      </div>
    </div>
  </div>
  </div>
</template>

<script setup>
const { fetchRecipe } = useRecipe()
const searchInput = ref('')
const result = ref(null)
const loading = ref(false)
const error = ref(null)

const searchRecipe = async () => {
  if (!searchInput.value.trim()) return

  loading.value = true
  error.value = null
  result.value = null 
  
  try {
    console.log('fetching recipe:', searchInput.value.trim())
    const data = await fetchRecipe(searchInput.value.trim())
    console.log('recipe data:', data) 
    console.log('recipe result:', data.recipe)
    result.value = data.recipe
    if (!result.value) {
      throw new Error('No recipe found')
    }
    console.log('recipe fetched:', result.steps)
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
    const path = window.location.pathname.slice(1) 
    
    if (path && path !== '') {
        searchInput.value = path
        searchRecipe()
        return
    }
    
    const subdomain = window.location.hostname.split('.')[0]
    if (subdomain && subdomain !== 'www' && subdomain !== 'localhost') {
        searchInput.value = subdomain
        searchRecipe()
    }})
</script>

<style scoped>
.container {
  min-height: 100vh;
  display: flex;
  padding: 1rem;
  font-family: system-ui, -apple-system, sans-serif;
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

.search-btn {
  background-color: #3b82f6;
  color: white;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s, transform 0.2s;
}

@media (max-width: 640px) {
  .container {
    margin: 1rem;
    padding: 0.5rem;
  }
  
  h1 {
    font-size: 1.5rem;
  }
}
</style>

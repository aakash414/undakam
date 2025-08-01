  export const useRecipe = () => {
    const config = useRuntimeConfig()
    
    const fetchRecipe = async (subdomain: string) => {
      try {
        const baseUrl = config.public.apiBase || 'http://localhost:8080/api'
        const url = `${baseUrl}/api`
        console.log('Fetching recipe from:', url, 'for subdomain:', subdomain)    
        const data = await $fetch(url, {
          method: 'GET',
          headers: { 
              'X-Forwarded-Host': subdomain,
              'Content-Type': 'application/json'
          }
        })
        console.log('Fetched data:', data)
        return data
      } catch (error) {
        console.error('Recipe fetch failed:', error)
        return null
      }
    }
    
    return { fetchRecipe }
  }
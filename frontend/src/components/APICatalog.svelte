<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { GetNetVisorAPIs, GetNetVisorAPIsByService, SearchNetVisorAPIs } from '../../wailsjs/wailsjs/go/main/App.js';

  const dispatch = createEventDispatcher();

  export let apiCatalog = null;

  let searchQuery = '';
  let selectedService = 'all';
  let selectedMethod = 'all';
  let selectedCategory = 'all';
  let netvisorAPIs = [];
  let isLoading = true;
  let hoveredAPI = null;

  onMount(async () => {
    console.log('APICatalog mounted, loading NetVisor APIs...');
    await loadNetVisorAPIs();
  });

  async function loadNetVisorAPIs() {
    try {
      console.log('Loading APIs from NetVisor database...');
      netvisorAPIs = await GetNetVisorAPIs();
      console.log('NetVisor APIs loaded:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to load NetVisor APIs:', error);
      netvisorAPIs = [];
    } finally {
      isLoading = false;
    }
  }

  async function searchAPIs() {
    if (!searchQuery.trim()) {
      await loadNetVisorAPIs();
      return;
    }
    
    try {
      console.log('Searching NetVisor APIs for:', searchQuery);
      netvisorAPIs = await SearchNetVisorAPIs(searchQuery);
      console.log('Search results:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to search NetVisor APIs:', error);
      netvisorAPIs = [];
    }
  }

  async function filterByService(service) {
    selectedService = service;
    try {
      if (service && service !== 'all') {
        console.log('Filtering NetVisor APIs by service:', service);
        netvisorAPIs = await GetNetVisorAPIsByService(service);
      } else {
        await loadNetVisorAPIs();
      }
      console.log('Filtered APIs:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to filter NetVisor APIs:', error);
      netvisorAPIs = [];
    }
  }

  function clearFilters() {
    searchQuery = '';
    selectedService = 'all';
    selectedMethod = 'all';
    selectedCategory = 'all';
    loadNetVisorAPIs();
  }

  function selectAPI(api) {
    dispatch('apiSelected', api);
  }

  // Get all unique categories from NetVisor APIs
  $: categories = netvisorAPIs
    .map(api => api.category)
    .filter((category, index, arr) => category && arr.indexOf(category) === index)
    .sort();

  // Get all unique methods from NetVisor APIs
  $: methods = netvisorAPIs
    .map(api => api.method)
    .filter((method, index, arr) => method && arr.indexOf(method) === index)
    .sort();

  // Get all unique services from NetVisor APIs
  $: services = netvisorAPIs
    .map(api => api.service)
    .filter((service, index, arr) => service && arr.indexOf(service) === index)
    .sort();

  // Filter APIs based on search and filters
  $: filteredAPIs = netvisorAPIs.filter(api => {
    const matchesSearch = !searchQuery || 
      (api.description && api.description.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (api.path && api.path.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (api.category && api.category.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (api.tags && api.tags.toLowerCase().includes(searchQuery.toLowerCase()));

    const matchesService = selectedService === 'all' || api.service === selectedService;
    const matchesMethod = selectedMethod === 'all' || api.method === selectedMethod;
    const matchesCategory = selectedCategory === 'all' || api.category === selectedCategory;

    return matchesSearch && matchesService && matchesMethod && matchesCategory;
  });

  function getMethodColor(method) {
    switch (method?.toUpperCase()) {
      case 'GET': return '#4CAF50';
      case 'POST': return '#2196F3';
      case 'PUT': return '#FF9800';
      case 'DELETE': return '#F44336';
      case 'RUNCMDS': return '#9C27B0';
      default: return '#9E9E9E';
    }
  }

  function getServiceColor(service) {
    switch (service?.toLowerCase()) {
      case 'eapi': return '#2196F3';
      case 'cloudvision': return '#9C27B0';
      case 'eos_rest': return '#FF9800';
      case 'telemetry': return '#4CAF50';
      default: return '#9E9E9E';
    }
  }
</script>

<div class="api-catalog">
  <div class="section-header">
    <h2>API Catalog</h2>
    <div class="api-count">
      {filteredAPIs.length} APIs
    </div>
  </div>

  <div class="search-filters">
    <div class="search-section">
      <input 
        type="text" 
        bind:value={searchQuery} 
        placeholder="Search APIs..." 
        class="search-input"
        on:input={searchAPIs}
      />
    </div>
    
    <div class="filter-section">
      <select bind:value={selectedService} on:change={() => filterByService(selectedService)} class="filter-select">
        <option value="all">All Services</option>
        {#each services as service}
          <option value={service}>{service.toUpperCase()}</option>
        {/each}
      </select>
      
      <select bind:value={selectedMethod} class="filter-select">
        <option value="all">All Methods</option>
        {#each methods as method}
          <option value={method}>{method.toUpperCase()}</option>
        {/each}
      </select>
      
      <select bind:value={selectedCategory} class="filter-select">
        <option value="all">All Categories</option>
        {#each categories as category}
          <option value={category}>{category}</option>
        {/each}
      </select>
      
      <button class="clear-filters-btn" on:click={clearFilters}>
        Clear Filters
      </button>
    </div>
  </div>

  {#if isLoading}
    <div class="loading">
      <p>Loading API catalog from NetVisor database...</p>
    </div>
  {:else if filteredAPIs.length === 0}
    <div class="no-apis">
      <p>No APIs found matching your criteria.</p>
      <p>Try adjusting your search terms or filters.</p>
    </div>
  {:else}
    <div class="api-list">
      {#each filteredAPIs as api}
        <div
          class="api-item"
          class:hovered={hoveredAPI?.id === api.id}
          on:click={() => selectAPI(api)}
          on:mouseenter={() => hoveredAPI = api}
          on:mouseleave={() => hoveredAPI = null}
        >
          <div class="api-header">
            <div class="api-method-service">
              <span class="method-badge" style="background-color: {getMethodColor(api.method)};">
                {api.method || 'GET'}
              </span>
              <span class="service-badge" style="background-color: {getServiceColor(api.service)};">
                {api.service || 'Unknown'}
              </span>
            </div>
            <div class="api-id">
              {api.id || 'Unknown ID'}
            </div>
          </div>
          
          <div class="api-content">
            <div class="api-path">
              {api.path || 'No path specified'}
            </div>
            <div class="api-description">
              {api.description || 'No description available'}
            </div>
            {#if api.category}
              <div class="api-category">
                <span class="category-label">Category:</span>
                <span class="category-value">{api.category}</span>
              </div>
            {/if}
            {#if api.tags}
              <div class="api-tags">
                <span class="tags-label">Tags:</span>
                <span class="tags-value">{api.tags}</span>
              </div>
            {/if}
            {#if api.parameters}
              <div class="api-parameters">
                <span class="parameters-label">Parameters:</span>
                <span class="parameters-value">{api.parameters}</span>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .api-catalog {
    max-width: 1400px;
    margin: 0 auto;
    padding: 1rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 0, 128, 0.2);
  }

  .section-header h2 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0;
  }

  .api-count {
    background: rgba(255, 0, 128, 0.1);
    border: 1px solid #ff0080;
    color: #ff0080;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-weight: bold;
  }

  .search-filters {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .search-section {
    display: flex;
    justify-content: center;
  }

  .search-input {
    width: 100%;
    max-width: 500px;
    padding: 0.75rem;
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
    color: white;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-size: 1rem;
  }

  .search-input:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .filter-section {
    display: flex;
    gap: 1rem;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
  }

  .filter-select {
    padding: 0.5rem;
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
    color: white;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    min-width: 120px;
  }

  .filter-select:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 5px rgba(255, 0, 128, 0.3);
  }

  .clear-filters-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
  }

  .clear-filters-btn:hover {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .loading, .no-apis {
    text-align: center;
    padding: 3rem;
    color: #e0e0e0;
  }

  .api-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    gap: 1rem;
  }

  .api-item {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .api-item:hover,
  .api-item.hovered {
    border-color: rgba(255, 0, 128, 0.5);
    background: rgba(255, 255, 255, 0.08);
    box-shadow: 0 0 15px rgba(255, 0, 128, 0.2);
    transform: translateY(-2px);
  }

  .api-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .api-method-service {
    display: flex;
    gap: 0.5rem;
  }

  .method-badge, .service-badge {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    color: white;
    font-size: 0.8rem;
    font-weight: bold;
    text-transform: uppercase;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .api-id {
    color: #888;
    font-size: 0.8rem;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .api-content {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .api-path {
    color: #ff0080;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
    font-size: 1rem;
    word-break: break-all;
  }

  .api-description {
    color: #e0e0e0;
    font-size: 0.9rem;
    line-height: 1.4;
  }

  .api-category, .api-tags, .api-parameters {
    display: flex;
    gap: 0.5rem;
    font-size: 0.8rem;
  }

  .category-label, .tags-label, .parameters-label {
    color: #888;
    font-weight: bold;
    min-width: 80px;
  }

  .category-value, .tags-value, .parameters-value {
    color: #ccc;
    flex: 1;
  }

  .category-value {
    color: #ff0080;
  }

  .tags-value {
    color: #4CAF50;
  }

  .parameters-value {
    color: #FF9800;
  }
</style>
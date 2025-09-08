<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { GetEndpoints, AddEndpoint, DeleteEndpoint, TestConnection } from '../../wailsjs/wailsjs/go/main/App.js';
  
  const dispatch = createEventDispatcher();
  
  let endpoints = [];
  let showAddForm = false;
  let newEndpoint = {
    name: '',
    type: 'eapi',
    url: '',
    username: '',
    password: '',
    token: ''
  };
  
  onMount(async () => {
    console.log('EndpointManager mounted');
    await loadEndpoints();
  });
  
  async function loadEndpoints() {
    try {
      console.log('Loading endpoints from database...');
      endpoints = await GetEndpoints();
      console.log('Endpoints loaded:', endpoints.length, 'endpoints');
    } catch (error) {
      console.error('Failed to load endpoints:', error);
      endpoints = [];
    }
  }
  
  async function addEndpoint() {
    console.log('Add endpoint called');
    if (newEndpoint.name && newEndpoint.url) {
      try {
        const endpoint = {
          name: newEndpoint.name,
          type: newEndpoint.type,
          url: newEndpoint.url,
          username: newEndpoint.username,
          password: newEndpoint.password,
          token: newEndpoint.token
        };
        
        await AddEndpoint(endpoint);
        
        // Reset form
        newEndpoint = {
          name: '',
          type: 'eapi',
          url: '',
          username: '',
          password: '',
          token: ''
        };
        showAddForm = false;
        
        // Reload endpoints from database
        await loadEndpoints();
        
        console.log('Endpoint added to database');
      } catch (error) {
        console.error('Error adding endpoint:', error);
        alert('Failed to add endpoint: ' + error.message);
      }
    }
  }
  
  function selectEndpoint(endpoint) {
    console.log('Endpoint selected:', endpoint);
    dispatch('endpointSelected', endpoint);
  }
  
  async function deleteEndpoint(id) {
    console.log('Delete endpoint called:', id);
    try {
      await DeleteEndpoint(id);
      await loadEndpoints();
      console.log('Endpoint deleted from database');
    } catch (error) {
      console.error('Error deleting endpoint:', error);
      alert('Failed to delete endpoint: ' + error.message);
    }
  }
  
  async function testConnection(endpoint) {
    console.log('Test connection called:', endpoint);
    try {
      // Set status to testing
      endpoint.status = 'Testing...';
      endpoints = [...endpoints];
      
      const result = await TestConnection(endpoint.id);
      endpoint.status = result.success ? 'Connected' : 'Failed';
      endpoints = [...endpoints];
      
      console.log('Connection test result:', result);
    } catch (error) {
      endpoint.status = 'Failed';
      endpoints = [...endpoints];
      console.error('Error testing connection:', error);
    }
  }
</script>

<div class="endpoint-manager">
  <div class="section-header">
    <h2>Endpoint Management</h2>
    <button class="action-btn" on:click={() => showAddForm = !showAddForm}>
      {showAddForm ? 'Cancel' : 'Add Endpoint'}
    </button>
  </div>

  {#if showAddForm}
    <div class="add-endpoint-form">
      <h3>Add New Endpoint</h3>
      <div class="form-grid">
        <div class="form-group">
          <label>Name</label>
          <input type="text" bind:value={newEndpoint.name} placeholder="Endpoint name" />
        </div>
        <div class="form-group">
          <label>Type</label>
          <select bind:value={newEndpoint.type}>
            <option value="eapi">EOS eAPI</option>
            <option value="cloudvision">CloudVision Portal</option>
            <option value="eos_rest">EOS REST API</option>
            <option value="telemetry">Telemetry</option>
          </select>
        </div>
        <div class="form-group">
          <label>URL</label>
          <input type="text" bind:value={newEndpoint.url} placeholder="https://..." />
        </div>
        {#if newEndpoint.type === 'eapi' || newEndpoint.type === 'eos_rest'}
          <div class="form-group">
            <label>Username</label>
            <input type="text" bind:value={newEndpoint.username} placeholder="Username" />
          </div>
          <div class="form-group">
            <label>Password</label>
            <input type="password" bind:value={newEndpoint.password} placeholder="Password" />
          </div>
        {:else if newEndpoint.type === 'cloudvision'}
          <div class="form-group">
            <label>Token</label>
            <input type="password" bind:value={newEndpoint.token} placeholder="Bearer token" />
          </div>
        {/if}
      </div>
      <div class="form-actions">
        <button class="save-btn" on:click={addEndpoint}>Save</button>
        <button class="cancel-btn" on:click={() => showAddForm = false}>Cancel</button>
      </div>
    </div>
  {/if}

  <div class="endpoints-list">
    {#each endpoints as endpoint}
      <div class="endpoint-card" on:click={() => selectEndpoint(endpoint)}>
        <div class="endpoint-header">
          <div class="endpoint-info">
            <h4>{endpoint.name}</h4>
            <span class="endpoint-type">{endpoint.type}</span>
          </div>
          <div class="endpoint-actions">
            <button class="test-btn" on:click|stopPropagation={() => testConnection(endpoint)}>
              Test
            </button>
            <button class="delete-btn" on:click|stopPropagation={() => deleteEndpoint(endpoint.id)}>
              Delete
            </button>
          </div>
        </div>
        <div class="endpoint-details">
          <p><strong>URL:</strong> {endpoint.url}</p>
          {#if endpoint.username}
            <p><strong>Username:</strong> {endpoint.username}</p>
          {/if}
          <p><strong>Status:</strong> <span class="status-{endpoint.status?.toLowerCase() || 'disconnected'}">{endpoint.status || 'Disconnected'}</span></p>
        </div>
      </div>
    {/each}
    
    {#if endpoints.length === 0}
      <div class="no-endpoints">
        <p>No endpoints configured. Click "Add Endpoint" to get started.</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .endpoint-manager {
    max-width: 1200px;
    margin: 0 auto;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .section-header h2 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0;
  }

  .action-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.3s ease;
  }

  .action-btn:hover {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .add-endpoint-form {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
    padding: 2rem;
    margin-bottom: 2rem;
  }

  .add-endpoint-form h3 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0 0 1.5rem 0;
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
  }

  .form-group label {
    color: #ccc;
    margin-bottom: 0.5rem;
    font-size: 0.9rem;
  }

  .form-group input,
  .form-group select {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: white;
    padding: 0.75rem;
    border-radius: 4px;
    font-size: 0.9rem;
  }

  .form-group input:focus,
  .form-group select:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 5px rgba(255, 0, 128, 0.3);
  }

  .form-actions {
    display: flex;
    gap: 1rem;
  }

  .save-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
  }

  .cancel-btn {
    background: transparent;
    border: 1px solid rgba(255, 255, 255, 0.3);
    color: #ccc;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
  }

  .endpoints-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .endpoint-card {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .endpoint-card:hover {
    border-color: rgba(255, 0, 128, 0.3);
    background: rgba(255, 255, 255, 0.08);
  }

  .endpoint-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .endpoint-info h4 {
    color: #ff0080;
    margin: 0 0 0.25rem 0;
    font-size: 1.1rem;
  }

  .endpoint-type {
    color: #888;
    font-size: 0.8rem;
    text-transform: uppercase;
  }

  .endpoint-actions {
    display: flex;
    gap: 0.5rem;
  }

  .test-btn {
    background: rgba(0, 255, 0, 0.2);
    border: 1px solid #00ff00;
    color: #00ff00;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .delete-btn {
    background: rgba(255, 0, 0, 0.2);
    border: 1px solid #ff0000;
    color: #ff0000;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .endpoint-details p {
    margin: 0.25rem 0;
    color: #ccc;
    font-size: 0.9rem;
  }

  .status-connected {
    color: #00ff00;
  }

  .status-disconnected {
    color: #ff0000;
  }

  .status-testing {
    color: #ffaa00;
  }

  .no-endpoints {
    text-align: center;
    padding: 3rem;
    color: #888;
  }
</style>
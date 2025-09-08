<script>
  import { createEventDispatcher, onMount } from 'svelte';

  const dispatch = createEventDispatcher();

  let apiConfigs = {
    eapi: {
      name: 'EOS eAPI',
      url: '',
      username: '',
      password: '',
      enabled: false,
      status: 'disconnected'
    },
    cloudvision: {
      name: 'CloudVision Portal',
      url: '',
      token: '',
      enabled: false,
      status: 'disconnected'
    },
    eos_rest: {
      name: 'EOS REST API',
      url: '',
      username: '',
      password: '',
      enabled: false,
      status: 'disconnected'
    }
  };

  let testResults = {};

  onMount(() => {
    // Load saved configurations
    loadConfigurations();
  });

  function loadConfigurations() {
    const saved = localStorage.getItem('arista_api_configs');
    if (saved) {
      try {
        apiConfigs = { ...apiConfigs, ...JSON.parse(saved) };
      } catch (error) {
        console.error('Failed to load API configurations:', error);
      }
    }
  }

  function saveConfigurations() {
    localStorage.setItem('arista_api_configs', JSON.stringify(apiConfigs));
  }

  async function testConnection(apiType) {
    const config = apiConfigs[apiType];
    if (!config.url) {
      testResults[apiType] = { success: false, message: 'URL is required' };
      return;
    }

    testResults[apiType] = { success: false, message: 'Testing...' };
    testResults = { ...testResults };

    try {
      // Call real backend API test
      const response = await fetch('/api/test-api-connection', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          type: apiType,
          url: config.url,
          username: config.username,
          password: config.password,
          token: config.token
        })
      });
      
      if (response.ok) {
        const result = await response.json();
        testResults[apiType] = { 
          success: result.success, 
          message: result.message,
          details: result.details || `Tested ${config.name} at ${config.url}`
        };
        config.status = result.success ? 'connected' : 'disconnected';
      } else {
        testResults[apiType] = { 
          success: false, 
          message: 'Connection test failed',
          details: 'Backend API test endpoint not available'
        };
        config.status = 'disconnected';
      }
    } catch (error) {
      testResults[apiType] = { 
        success: false, 
        message: 'Connection failed',
        details: error.message
      };
      config.status = 'disconnected';
    }

    testResults = { ...testResults };
    saveConfigurations();
  }

  function updateConfig(apiType, field, value) {
    apiConfigs[apiType][field] = value;
    saveConfigurations();
  }

  function toggleAPI(apiType) {
    apiConfigs[apiType].enabled = !apiConfigs[apiType].enabled;
    saveConfigurations();
  }

  function getStatusColor(status) {
    switch (status) {
      case 'connected': return '#00ff88';
      case 'disconnected': return '#ff0080';
      default: return '#888';
    }
  }

  function getTestResultColor(apiType) {
    const result = testResults[apiType];
    if (!result) return '#888';
    return result.success ? '#00ff88' : '#ff0080';
  }
</script>

<div class="api-manager">
  <div class="header">
    <h2>API Connection Manager</h2>
    <p>Configure and test connections to all Arista API endpoints</p>
  </div>

  <div class="api-grid">
    {#each Object.entries(apiConfigs) as [apiType, config]}
      <div class="api-card">
        <div class="api-header">
          <div class="api-title">
            <h3>{config.name}</h3>
            <div class="api-status">
              <div class="status-dot" style="background-color: {getStatusColor(config.status)}"></div>
              <span class="status-text">{config.status}</span>
            </div>
          </div>
          <label class="toggle-switch">
            <input 
              type="checkbox" 
              checked={config.enabled}
              on:change={() => toggleAPI(apiType)}
            />
            <span class="toggle-slider"></span>
          </label>
        </div>

        <div class="api-config">
          <div class="form-group">
            <label for="{apiType}-url">API URL:</label>
            <input 
              id="{apiType}-url"
              type="url" 
              bind:value={config.url}
              on:input={(e) => updateConfig(apiType, 'url', e.target.value)}
              placeholder="https://device-ip or https://cvp-server"
            />
          </div>

          {#if apiType === 'cloudvision'}
            <div class="form-group">
              <label for="{apiType}-token">API Token:</label>
              <input 
                id="{apiType}-token"
                type="password" 
                bind:value={config.token}
                on:input={(e) => updateConfig(apiType, 'token', e.target.value)}
                placeholder="Bearer token for CloudVision"
              />
            </div>
          {:else}
            <div class="form-group">
              <label for="{apiType}-username">Username:</label>
              <input 
                id="{apiType}-username"
                type="text" 
                bind:value={config.username}
                on:input={(e) => updateConfig(apiType, 'username', e.target.value)}
                placeholder="API username"
              />
            </div>
            <div class="form-group">
              <label for="{apiType}-password">Password:</label>
              <input 
                id="{apiType}-password"
                type="password" 
                bind:value={config.password}
                on:input={(e) => updateConfig(apiType, 'password', e.target.value)}
                placeholder="API password"
              />
            </div>
          {/if}

          <div class="test-section">
            <button 
              class="test-btn"
              on:click={() => testConnection(apiType)}
              disabled={!config.url}
            >
              Test Connection
            </button>
            
            {#if testResults[apiType]}
              <div class="test-result" style="color: {getTestResultColor(apiType)}">
                <strong>{testResults[apiType].message}</strong>
                {#if testResults[apiType].details}
                  <div class="test-details">{testResults[apiType].details}</div>
                {/if}
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>

  <div class="summary">
    <h3>Connection Summary</h3>
    <div class="summary-stats">
      {#each Object.entries(apiConfigs) as [apiType, config]}
        <div class="stat-item">
          <span class="stat-label">{config.name}:</span>
          <span class="stat-value" style="color: {getStatusColor(config.status)}">
            {config.enabled ? config.status : 'disabled'}
          </span>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .api-manager {
    padding: 20px;
    color: white;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .header {
    margin-bottom: 30px;
    text-align: center;
  }

  .header h2 {
    font-size: 1.8rem;
    margin: 0 0 10px 0;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    background: linear-gradient(45deg, #ff0080, #ff4080);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 0 30px rgba(255, 0, 128, 0.5);
  }

  .header p {
    color: #888;
    margin: 0;
  }

  .api-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 20px;
    margin-bottom: 30px;
  }

  .api-card {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    border-radius: 8px;
    padding: 20px;
    transition: all 0.3s ease;
  }

  .api-card:hover {
    border-color: #ff4080;
    box-shadow: 0 0 20px rgba(255, 0, 128, 0.2);
  }

  .api-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .api-title h3 {
    margin: 0 0 5px 0;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    color: #ff4080;
  }

  .api-status {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    transition: all 0.3s ease;
  }

  .status-text {
    font-size: 0.9rem;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
  }

  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 50px;
    height: 24px;
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #333;
    transition: 0.3s;
    border-radius: 24px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  input:checked + .toggle-slider {
    background-color: #ff0080;
  }

  input:checked + .toggle-slider:before {
    transform: translateX(26px);
  }

  .form-group {
    margin-bottom: 15px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    color: #ff4080;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
  }

  .form-group input {
    width: 100%;
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid #444;
    border-radius: 4px;
    padding: 10px;
    color: white;
    font-family: inherit;
  }

  .form-group input:focus {
    outline: none;
    border-color: #ff4080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .test-section {
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid #333;
  }

  .test-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: none;
    color: white;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
    transition: all 0.3s ease;
  }

  .test-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 15px rgba(255, 0, 128, 0.4);
  }

  .test-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .test-result {
    margin-top: 10px;
    padding: 10px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 4px;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .test-details {
    font-size: 0.9rem;
    margin-top: 5px;
    opacity: 0.8;
  }

  .summary {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    border-radius: 8px;
    padding: 20px;
  }

  .summary h3 {
    margin: 0 0 15px 0;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    color: #ff4080;
  }

  .summary-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 15px;
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-label {
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
  }

  .stat-value {
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
  }
</style>

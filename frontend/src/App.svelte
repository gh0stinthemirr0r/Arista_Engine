<script>
  import { onMount } from 'svelte';
  import EndpointManager from './components/EndpointManager.svelte';
  import APIExplorer from './components/APIExplorer.svelte';
  import ResponseViewer from './components/ResponseViewer.svelte';
  import APICatalog from './components/APICatalog.svelte';
  import QueryLog from './components/QueryLog.svelte';
  import CloudVisionManager from './components/CloudVisionManager.svelte';
  import APIManager from './components/APIManager.svelte';
  import DeviceInventory from './components/DeviceInventory.svelte';
  import { GetAPICatalog } from '../wailsjs/wailsjs/go/main/App.js';
  
  let activeTab = 'start';
  let selectedEndpoint = null;
  let apiResponse = null;
  let apiCatalog = { eapi: {}, cloudvision: {}, eos_rest: {}, telemetry: {} };
  let queryLog = [];
  let showGettingStarted = true;
  
  onMount(async () => {
    console.log('App onMount triggered.');
    try {
      // Load API catalog from backend
      console.log('Loading API catalog from backend...');
      apiCatalog = await GetAPICatalog();
      console.log('API catalog loaded:', Object.keys(apiCatalog));
    } catch (error) {
      console.error('Failed to load API catalog:', error);
      // Fallback to empty catalog
      apiCatalog = { eapi: {}, cloudvision: {}, eos_rest: {}, telemetry: {} };
    }
    
    // Initialize query log
    queryLog = [];
    console.log('App initialization complete.');
  });

  function handleEndpointSelected(event) {
    console.log('Endpoint selected:', event.detail);
    selectedEndpoint = event.detail;
  }

  function handleAPIResponse(response) {
    console.log('API response received:', response);
    apiResponse = response;
    activeTab = 'response';
  }
  
  function handleTabChange(tab) {
    console.log('Tab change requested:', tab);
    activeTab = tab;
    if (tab !== 'start') {
      showGettingStarted = false;
    }
    console.log('Active tab set to:', activeTab);
  }

  function startGettingStarted() {
    console.log('Starting getting started flow');
    showGettingStarted = true;
    activeTab = 'start';
  }

  function proceedToEndpoints() {
    console.log('Proceeding to endpoints');
    showGettingStarted = false;
    activeTab = 'endpoints';
  }
</script>

<main>
  <div class="app-container">
    <header class="app-header">
      <div class="header-content">
        <h1>Arista_Engine</h1>
      </div>
      <div class="header-status">
        <div class="status-indicator">
          <div class="status-dot" class:active={selectedEndpoint !== null}></div>
          <span class="status-text">{selectedEndpoint ? 'Online' : 'Offline'}</span>
        </div>
      </div>
    </header>

    <nav class="tab-navigation">
      <button
        class="tab-button start-button"
        class:active={activeTab === 'start'}
        on:click={startGettingStarted}
      >
        Start
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'endpoints'}
        on:click={() => handleTabChange('endpoints')}
      >
        Endpoints
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'catalog'}
        on:click={() => handleTabChange('catalog')}
      >
        API Catalog
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'explorer'}
        on:click={() => handleTabChange('explorer')}
        disabled={!selectedEndpoint}
      >
        Explorer
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'cloudvision'}
        on:click={() => handleTabChange('cloudvision')}
        disabled={!selectedEndpoint || selectedEndpoint?.type !== 'cloudvision'}
      >
        CloudVision
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'response'}
        on:click={() => handleTabChange('response')}
        disabled={!apiResponse}
      >
        Response
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'logs'}
        on:click={() => handleTabChange('logs')}
      >
        Logs
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'inventory'}
        on:click={() => handleTabChange('inventory')}
      >
        Inventory
      </button>
      <button
        class="tab-button"
        class:active={activeTab === 'api'}
        on:click={() => handleTabChange('api')}
      >
        API
      </button>
    </nav>

    <div class="main-content">
      {#if activeTab === 'start'}
        <div class="getting-started">
          <div class="getting-started-content">
            <h2>Welcome to Arista Engine</h2>
            <p class="intro-text">
              Arista Engine is your comprehensive tool for managing and exploring Arista Networks devices and APIs.
              Get started by adding your first device endpoint to unlock the full power of the application.
            </p>

            <div class="flow-steps">
              <div class="step">
                <div class="step-number">1</div>
                <div class="step-content">
                  <h3>Add Device Endpoint</h3>
                  <p>Navigate to the Endpoints tab and add your Arista device (EOS eAPI, CloudVision, or EOS REST).
                     Provide the device URL, credentials, and connection details.</p>
                </div>
              </div>

              <div class="step">
                <div class="step-number">2</div>
                <div class="step-content">
                  <h3>Explore APIs</h3>
                  <p>Once connected, browse the comprehensive API catalog with thousands of available commands
                     and endpoints across all Arista services.</p>
                </div>
              </div>

              <div class="step">
                <div class="step-number">3</div>
                <div class="step-content">
                  <h3>Execute Commands</h3>
                  <p>Use the API Explorer to run commands, manage CloudVision Portal, and interact with your devices
                     through a powerful, user-friendly interface.</p>
                </div>
              </div>
            </div>

            <div class="features-preview">
              <h3>Available Features</h3>
              <div class="feature-grid">
                <div class="feature-item">
                  <strong>EOS eAPI</strong>
                  <span>Execute CLI commands via JSON-RPC</span>
                </div>
                <div class="feature-item">
                  <strong>CloudVision Portal</strong>
                  <span>Manage devices, configlets, and tasks</span>
                </div>
                <div class="feature-item">
                  <strong>EOS REST API</strong>
                  <span>RESTful interface for EOS management</span>
                </div>
                <div class="feature-item">
                  <strong>Telemetry</strong>
                  <span>Stream and analyze device telemetry</span>
                </div>
              </div>
            </div>

            <div class="get-started-actions">
              <button class="get-started-button" on:click={proceedToEndpoints}>
                Get Started - Add Your First Device
              </button>
              <p class="note">
                Note: Most features are disabled until you add and select a device endpoint.
              </p>
            </div>
          </div>
        </div>
      {:else if activeTab === 'endpoints'}
        <EndpointManager on:endpointSelected={handleEndpointSelected} />
      {:else if activeTab === 'catalog'}
        <APICatalog {apiCatalog} />
      {:else if activeTab === 'explorer'}
        <APIExplorer {selectedEndpoint} {apiCatalog} on:response={handleAPIResponse} />
      {:else if activeTab === 'cloudvision'}
        <CloudVisionManager {selectedEndpoint} on:response={handleAPIResponse} />
      {:else if activeTab === 'response'}
        <ResponseViewer response={apiResponse} />
      {:else if activeTab === 'logs'}
        <QueryLog {queryLog} />
      {:else if activeTab === 'inventory'}
        <DeviceInventory />
      {:else if activeTab === 'api'}
        <APIManager />
      {/if}
    </div>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background: #0a0a0a;
    color: #ffffff;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    overflow: hidden;
  }

  .app-container {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
  }

  .app-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 2rem;
    background: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 0, 128, 0.3);
    backdrop-filter: blur(10px);
  }

  .header-content h1 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 1.8rem;
    font-weight: bold;
    color: #ff0080;
    margin: 0;
    text-shadow: 0 0 10px rgba(255, 0, 128, 0.5);
  }

  .header-status {
    display: flex;
    align-items: center;
  }

  .status-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .status-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #666;
    transition: all 0.3s ease;
  }

  .status-dot.active {
    background: #00ff00;
    box-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
  }

  .status-text {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 0.9rem;
    color: #ccc;
  }

  .tab-navigation {
    display: flex;
    background: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 0, 128, 0.2);
    overflow-x: auto;
    flex-wrap: nowrap;
    padding: 0.5rem;
  }

  .tab-button {
    background: transparent;
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #ccc;
    padding: 0.75rem 1.5rem;
    margin-right: 0.5rem;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 0.9rem;
    white-space: nowrap;
  }

  .tab-button:hover {
    background: rgba(255, 0, 128, 0.1);
    border-color: rgba(255, 0, 128, 0.3);
    color: #fff;
  }

  .tab-button.active {
    background: rgba(255, 0, 128, 0.2);
    border-color: #ff0080;
    color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .tab-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .start-button {
    background: transparent;
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #ff0080;
    font-weight: bold;
    padding: 0.75rem 1.5rem;
    min-width: 80px;
    white-space: nowrap;
    overflow: visible;
    flex-shrink: 0;
    margin-right: 0.5rem;
    border-radius: 4px;
  }

  .start-button:hover {
    background: rgba(255, 0, 128, 0.1);
    border-color: rgba(255, 0, 128, 0.3);
    color: #ff0080;
  }

  .start-button.active {
    background: rgba(255, 0, 128, 0.2);
    border-color: #ff0080;
    color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .main-content {
    flex: 1;
    padding: 2rem;
    overflow-y: auto;
    background: rgba(0, 0, 0, 0.3);
  }

  .getting-started {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }

  .getting-started-content {
    text-align: center;
  }

  .getting-started h2 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 2.5rem;
    color: #ff0080;
    margin-bottom: 1rem;
    text-shadow: 0 0 20px rgba(255, 0, 128, 0.5);
  }

  .intro-text {
    font-size: 1.1rem;
    color: #ccc;
    margin-bottom: 3rem;
    line-height: 1.6;
  }

  .flow-steps {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    margin-bottom: 3rem;
  }

  .step {
    display: flex;
    align-items: flex-start;
    gap: 1.5rem;
    text-align: left;
  }

  .step-number {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: linear-gradient(45deg, #ff0080, #ff4080);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 1.2rem;
    flex-shrink: 0;
  }

  .step-content h3 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0 0 0.5rem 0;
    font-size: 1.3rem;
  }

  .step-content p {
    color: #ccc;
    margin: 0;
    line-height: 1.5;
  }

  .features-preview {
    margin-bottom: 3rem;
  }

  .features-preview h3 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin-bottom: 1.5rem;
    font-size: 1.5rem;
  }

  .feature-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
  }

  .feature-item {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
    padding: 1.5rem;
    text-align: left;
  }

  .feature-item strong {
    display: block;
    color: #ff0080;
    font-size: 1.1rem;
    margin-bottom: 0.5rem;
  }

  .feature-item span {
    color: #ccc;
    font-size: 0.9rem;
  }

  .get-started-actions {
    text-align: center;
  }

  .get-started-button {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 1rem 2rem;
    border-radius: 8px;
    font-size: 1.1rem;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    box-shadow: 0 4px 15px rgba(255, 0, 128, 0.3);
  }

  .get-started-button:hover {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 6px 20px rgba(255, 0, 128, 0.5);
    transform: translateY(-2px);
  }

  .note {
    color: #888;
    font-size: 0.9rem;
    margin-top: 1rem;
    font-style: italic;
  }
</style>
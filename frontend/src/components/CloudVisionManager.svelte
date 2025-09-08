<script>
  import { createEventDispatcher, onMount } from 'svelte';

  const dispatch = createEventDispatcher();

  export let selectedEndpoint = null;

  let activeSection = 'inventory';
  let inventory = [];
  let configlets = [];
  let tasks = [];
  let isLoading = false;
  let response = null;

  // Form data for different operations
  let newConfiglet = {
    name: '',
    content: ''
  };
  let selectedDevice = '';
  let selectedConfiglet = '';
  let taskId = '';

  onMount(async () => {
    // Load real data from CloudVision Portal
    await loadInventory();
    await loadConfiglets();
    await loadTasks();
  });

  async function loadInventory() {
    try {
      // TODO: Implement real CloudVision API call
      inventory = [];
    } catch (error) {
      console.error('Failed to load inventory:', error);
      inventory = [];
    }
  }

  async function loadConfiglets() {
    try {
      // TODO: Implement real CloudVision API call
      configlets = [];
    } catch (error) {
      console.error('Failed to load configlets:', error);
      configlets = [];
    }
  }

  async function loadTasks() {
    try {
      // TODO: Implement real CloudVision API call
      tasks = [];
    } catch (error) {
      console.error('Failed to load tasks:', error);
      tasks = [];
    }
  }

  function getMethodColor(method) {
    const colors = {
      'GET': '#00ff88',
      'POST': '#ff8800',
      'PUT': '#0088ff',
      'DELETE': '#ff0088'
    };
    return colors[method] || '#ffffff';
  }

  async function getInventory() {
    isLoading = true;
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      response = {
        method: 'GET',
        path: '/api/v1/inventory',
        status: 200,
        data: inventory,
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
    } catch (error) {
      response = {
        method: 'GET',
        path: '/api/v1/inventory',
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }

  async function getConfiglets() {
    isLoading = true;
    try {
      await new Promise(resolve => setTimeout(resolve, 1000));
      response = {
        method: 'GET',
        path: '/api/v1/configlets',
        status: 200,
        data: configlets,
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
    } catch (error) {
      response = {
        method: 'GET',
        path: '/api/v1/configlets',
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }

  async function createConfiglet() {
    if (!newConfiglet.name || !newConfiglet.content) return;
    
    isLoading = true;
    try {
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const newConfigletData = {
        key: `configlet_${Date.now()}`,
        name: newConfiglet.name,
        config: newConfiglet.content,
        containerCount: 0,
        netElementCount: 0
      };
      
      configlets.push(newConfigletData);
      
      response = {
        method: 'POST',
        path: '/api/v1/configlets',
        status: 201,
        data: newConfigletData,
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
      
      // Reset form
      newConfiglet = { name: '', content: '' };
    } catch (error) {
      response = {
        method: 'POST',
        path: '/api/v1/configlets',
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }

  async function applyConfiglet() {
    if (!selectedDevice || !selectedConfiglet) return;
    
    isLoading = true;
    try {
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      const taskData = {
        id: Date.now().toString(),
        name: `Apply ${selectedConfiglet} to ${selectedDevice}`,
        status: 'Pending',
        created: new Date().toISOString()
      };
      
      tasks.push(taskData);
      
      response = {
        method: 'POST',
        path: '/api/v1/tasks',
        status: 201,
        data: { status: 'success', taskIds: [taskData.id] },
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
    } catch (error) {
      response = {
        method: 'POST',
        path: '/api/v1/tasks',
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }

  async function executeTask() {
    if (!taskId) return;
    
    isLoading = true;
    try {
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Update task status
      const task = tasks.find(t => t.id === taskId);
      if (task) {
        task.status = 'Completed';
      }
      
      response = {
        method: 'POST',
        path: `/api/v1/tasks/${taskId}/execute`,
        status: 200,
        data: { status: 'success', message: 'Task executed successfully' },
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
    } catch (error) {
      response = {
        method: 'POST',
        path: `/api/v1/tasks/${taskId}/execute`,
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }

  async function getCvpInfo() {
    isLoading = true;
    try {
      await new Promise(resolve => setTimeout(resolve, 500));
      response = {
        method: 'GET',
        path: '/cvpInfo/getCvpInfo.do',
        status: 200,
        data: { version: 'cvp_version', appVersion: 'app_version' },
        headers: { 'Content-Type': 'application/json' }
      };
      dispatch('response', response);
    } catch (error) {
      response = {
        method: 'GET',
        path: '/cvpInfo/getCvpInfo.do',
        status: 500,
        data: { error: error.message },
        headers: {}
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="cloudvision-manager">
  <div class="header">
    <h2>CloudVision Portal Manager</h2>
    <p>Manage devices, configlets, and tasks in CloudVision Portal</p>
  </div>

  {#if !selectedEndpoint}
    <div class="no-endpoint">
      <p>Please select a CloudVision endpoint first</p>
    </div>
  {:else}
    <div class="cloudvision-content">
      <!-- Navigation Tabs -->
      <div class="nav-tabs">
        <button 
          class="tab" 
          class:active={activeSection === 'inventory'}
          on:click={() => activeSection = 'inventory'}
        >
          Device Inventory
        </button>
        <button 
          class="tab" 
          class:active={activeSection === 'configlets'}
          on:click={() => activeSection = 'configlets'}
        >
          Configlets
        </button>
        <button 
          class="tab" 
          class:active={activeSection === 'tasks'}
          on:click={() => activeSection = 'tasks'}
        >
          Tasks
        </button>
        <button 
          class="tab" 
          class:active={activeSection === 'info'}
          on:click={() => activeSection = 'info'}
        >
          CVP Info
        </button>
      </div>

      <!-- Device Inventory Section -->
      {#if activeSection === 'inventory'}
        <div class="section">
          <div class="section-header">
            <h3>Device Inventory</h3>
            <button class="action-btn" on:click={getInventory} disabled={isLoading}>
              {isLoading ? 'Loading...' : 'Refresh Inventory'}
            </button>
          </div>
          
          <div class="inventory-grid">
            {#each inventory as device}
              <div class="device-card">
                <div class="device-header">
                  <h4>{device.fqdn}</h4>
                  <span class="device-type">{device.type}</span>
                </div>
                <div class="device-details">
                  <p><strong>IP:</strong> {device.ipAddress}</p>
                  <p><strong>Model:</strong> {device.modelName}</p>
                  <p><strong>Version:</strong> {device.version}</p>
                  <p><strong>Serial:</strong> {device.serialNumber}</p>
                  <p><strong>MAC:</strong> {device.systemMacAddress}</p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Configlets Section -->
      {#if activeSection === 'configlets'}
        <div class="section">
          <div class="section-header">
            <h3>Configlet Management</h3>
            <button class="action-btn" on:click={getConfiglets} disabled={isLoading}>
              {isLoading ? 'Loading...' : 'Refresh Configlets'}
            </button>
          </div>

          <!-- Create New Configlet -->
          <div class="configlet-form">
            <h4>Create New Configlet</h4>
            <div class="form-group">
              <label for="configlet-name">Configlet Name:</label>
              <input 
                id="configlet-name"
                type="text" 
                bind:value={newConfiglet.name} 
                placeholder="e.g., VLANS"
              />
            </div>
            <div class="form-group">
              <label for="configlet-content">Configuration Content:</label>
              <textarea 
                id="configlet-content"
                bind:value={newConfiglet.content} 
                placeholder="vlan 100&#10;   name DEMO&#10;end"
                rows="6"
              ></textarea>
            </div>
            <button class="action-btn" on:click={createConfiglet} disabled={isLoading || !newConfiglet.name || !newConfiglet.content}>
              {isLoading ? 'Creating...' : 'Create Configlet'}
            </button>
          </div>

          <!-- Apply Configlet to Device -->
          <div class="apply-form">
            <h4>Apply Configlet to Device</h4>
            <div class="form-row">
              <div class="form-group">
                <label for="select-device">Device:</label>
                <select id="select-device" bind:value={selectedDevice}>
                  <option value="">Select a device...</option>
                  {#each inventory as device}
                    <option value={device.fqdn}>{device.fqdn} ({device.ipAddress})</option>
                  {/each}
                </select>
              </div>
              <div class="form-group">
                <label for="select-configlet">Configlet:</label>
                <select id="select-configlet" bind:value={selectedConfiglet}>
                  <option value="">Select a configlet...</option>
                  {#each configlets as configlet}
                    <option value={configlet.name}>{configlet.name}</option>
                  {/each}
                </select>
              </div>
            </div>
            <button class="action-btn" on:click={applyConfiglet} disabled={isLoading || !selectedDevice || !selectedConfiglet}>
              {isLoading ? 'Applying...' : 'Apply Configlet'}
            </button>
          </div>

          <!-- Configlets List -->
          <div class="configlets-list">
            <h4>Existing Configlets</h4>
            {#each configlets as configlet}
              <div class="configlet-card">
                <div class="configlet-header">
                  <h5>{configlet.name}</h5>
                  <span class="configlet-stats">
                    {configlet.netElementCount} device(s)
                  </span>
                </div>
                <div class="configlet-content">
                  <pre>{configlet.config}</pre>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Tasks Section -->
      {#if activeSection === 'tasks'}
        <div class="section">
          <div class="section-header">
            <h3>Task Management</h3>
          </div>

          <!-- Execute Task -->
          <div class="task-form">
            <h4>Execute Task</h4>
            <div class="form-group">
              <label for="task-id">Task ID:</label>
              <input 
                id="task-id"
                type="text" 
                bind:value={taskId} 
                placeholder="Enter task ID to execute"
              />
            </div>
            <button class="action-btn" on:click={executeTask} disabled={isLoading || !taskId}>
              {isLoading ? 'Executing...' : 'Execute Task'}
            </button>
          </div>

          <!-- Tasks List -->
          <div class="tasks-list">
            <h4>Recent Tasks</h4>
            {#each tasks as task}
              <div class="task-card">
                <div class="task-header">
                  <h5>{task.name}</h5>
                  <span class="task-status" class:completed={task.status === 'Completed'}>
                    {task.status}
                  </span>
                </div>
                <div class="task-details">
                  <p><strong>ID:</strong> {task.id}</p>
                  <p><strong>Created:</strong> {new Date(task.created).toLocaleString()}</p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- CVP Info Section -->
      {#if activeSection === 'info'}
        <div class="section">
          <div class="section-header">
            <h3>CloudVision Portal Information</h3>
            <button class="action-btn" on:click={getCvpInfo} disabled={isLoading}>
              {isLoading ? 'Loading...' : 'Get CVP Info'}
            </button>
          </div>
          
          <div class="info-content">
            <p>Get information about the CloudVision Portal appliance including version details.</p>
            <div class="info-card">
              <h4>Connection Details</h4>
              <p><strong>Endpoint:</strong> {selectedEndpoint?.name}</p>
              <p><strong>URL:</strong> {selectedEndpoint?.url}</p>
              <p><strong>Type:</strong> {selectedEndpoint?.type}</p>
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .cloudvision-manager {
    padding: 20px;
    color: white;
    font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .header {
    margin-bottom: 30px;
    text-align: center;
  }

  .header h2 {
    font-size: 1.8rem;
    margin: 0 0 10px 0;
    background: linear-gradient(45deg, #ff0080, #ff4080);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 0 30px rgba(255, 0, 128, 0.5);
  }

  .header p {
    color: #888;
    margin: 0;
  }

  .no-endpoint {
    text-align: center;
    padding: 40px;
    color: #888;
  }

  .nav-tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 30px;
    border-bottom: 2px solid #333;
  }

  .tab {
    background: transparent;
    border: none;
    color: #888;
    padding: 12px 20px;
    cursor: pointer;
    border-bottom: 2px solid transparent;
    transition: all 0.3s ease;
    font-family: inherit;
  }

  .tab:hover {
    color: #ff4080;
    border-bottom-color: #ff4080;
  }

  .tab.active {
    color: #ff0080;
    border-bottom-color: #ff0080;
    background: rgba(255, 0, 128, 0.1);
  }

  .section {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    border-radius: 8px;
    padding: 20px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .section-header h3 {
    margin: 0;
    color: #ff4080;
  }

  .action-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: none;
    color: white;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-family: inherit;
    font-weight: bold;
    transition: all 0.3s ease;
  }

  .action-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 15px rgba(255, 0, 128, 0.4);
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .inventory-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }

  .device-card {
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid #444;
    border-radius: 8px;
    padding: 15px;
    transition: all 0.3s ease;
  }

  .device-card:hover {
    border-color: #ff4080;
    box-shadow: 0 0 20px rgba(255, 0, 128, 0.2);
  }

  .device-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }

  .device-header h4 {
    margin: 0;
    color: #ff4080;
  }

  .device-type {
    background: #333;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
  }

  .device-details p {
    margin: 5px 0;
    font-size: 0.9rem;
  }

  .form-group {
    margin-bottom: 15px;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    color: #ff4080;
    font-weight: bold;
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid #444;
    border-radius: 4px;
    padding: 10px;
    color: white;
    font-family: inherit;
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #ff4080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .configlet-form,
  .apply-form,
  .task-form {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
  }

  .configlet-form h4,
  .apply-form h4,
  .task-form h4 {
    margin: 0 0 15px 0;
    color: #ff4080;
  }

  .configlets-list,
  .tasks-list {
    margin-top: 20px;
  }

  .configlets-list h4,
  .tasks-list h4 {
    margin: 0 0 15px 0;
    color: #ff4080;
  }

  .configlet-card,
  .task-card {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #444;
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 15px;
  }

  .configlet-header,
  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }

  .configlet-header h5,
  .task-header h5 {
    margin: 0;
    color: #ff4080;
  }

  .configlet-stats {
    background: #333;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
  }

  .task-status {
    background: #ff8800;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
  }

  .task-status.completed {
    background: #00ff88;
    color: #000;
  }

  .configlet-content pre {
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid #444;
    border-radius: 4px;
    padding: 10px;
    margin: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.9rem;
    color: #00ff88;
    overflow-x: auto;
  }

  .task-details p {
    margin: 5px 0;
    font-size: 0.9rem;
  }

  .info-content {
    text-align: center;
  }

  .info-content p {
    color: #888;
    margin-bottom: 20px;
  }

  .info-card {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    border-radius: 8px;
    padding: 20px;
    max-width: 400px;
    margin: 0 auto;
  }

  .info-card h4 {
    margin: 0 0 15px 0;
    color: #ff4080;
  }

  .info-card p {
    margin: 8px 0;
    text-align: left;
  }
</style>

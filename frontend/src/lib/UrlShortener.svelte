<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>ktrk.ee - Premium URL Shortener</title>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;600;700&display=swap">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
  <style>
    :root {
      --primary: #00ff41;
      --primary-dark: #008f24;
      --bg: #0a0a0a;
      --bg-light: #121212;
      --bg-card: #1a1a1a;
      --text: #f0f0f0;
      --text-secondary: #aaa;
      --error: #ff3860;
      --success: #23d160;
    }
    
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
      font-family: 'JetBrains Mono', monospace;
    }
    
    body {
      background: var(--bg);
      color: var(--text);
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 2rem 1rem;
      background-image: 
        radial-gradient(circle at 10% 20%, rgba(0, 255, 65, 0.05) 0%, transparent 20%),
        radial-gradient(circle at 90% 80%, rgba(0, 255, 65, 0.05) 0%, transparent 20%);
      position: relative;
      overflow-x: hidden;
    }
    
    body::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: linear-gradient(90deg, var(--primary), var(--primary-dark));
      box-shadow: 0 0 15px var(--primary);
    }
    
    .container {
      width: 100%;
      max-width: 600px;
      display: flex;
      flex-direction: column;
      gap: 2rem;
      z-index: 1;
    }
    
    header {
      text-align: center;
      padding: 1rem 0;
      border-bottom: 1px solid rgba(0, 255, 65, 0.2);
    }
    
    .logo {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.75rem;
      margin-bottom: 0.5rem;
    }
    
    .logo-icon {
      color: var(--primary);
      font-size: 2.5rem;
    }
    
    h1 {
      font-size: 2.5rem;
      font-weight: 700;
      background: linear-gradient(90deg, var(--primary), #00cc55);
      -webkit-background-clip: text;
      background-clip: text;
      color: transparent;
      letter-spacing: -1px;
    }
    
    .tagline {
      color: var(--text-secondary);
      font-size: 1.1rem;
      margin-top: 0.25rem;
    }
    
    .stats-container {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 1rem;
      margin-top: 1rem;
    }
    
    .stat {
      background: var(--bg-card);
      border: 1px solid rgba(0, 255, 65, 0.1);
      border-radius: 8px;
      padding: 1rem;
      text-align: center;
      transition: all 0.3s ease;
    }
    
    .stat:hover {
      border-color: var(--primary);
      transform: translateY(-3px);
      box-shadow: 0 5px 15px rgba(0, 255, 65, 0.1);
    }
    
    .stat-value {
      font-size: 2rem;
      font-weight: 700;
      color: var(--primary);
      margin-bottom: 0.25rem;
    }
    
    .stat-label {
      font-size: 0.9rem;
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 1px;
    }
    
    .card {
      background: var(--bg-card);
      border-radius: 12px;
      padding: 2rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
      border: 1px solid rgba(0, 255, 65, 0.1);
      transition: all 0.3s ease;
    }
    
    .card:hover {
      border-color: var(--primary);
      box-shadow: 0 10px 30px rgba(0, 255, 65, 0.15);
    }
    
    .card-title {
      font-size: 1.5rem;
      margin-bottom: 1.5rem;
      display: flex;
      align-items: center;
      gap: 0.75rem;
      color: var(--primary);
    }
    
    .card-title i {
      font-size: 1.25rem;
    }
    
    .form-group {
      margin-bottom: 1.5rem;
    }
    
    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
      color: var(--text-secondary);
      font-size: 0.9rem;
    }
    
    .input-wrapper {
      position: relative;
    }
    
    .input-wrapper i {
      position: absolute;
      left: 15px;
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-secondary);
    }
    
    input {
      width: 100%;
      padding: 0.9rem 1rem 0.9rem 3rem;
      background: rgba(20, 20, 20, 0.7);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 8px;
      color: var(--text);
      font-size: 1rem;
      transition: all 0.3s ease;
    }
    
    input:focus {
      outline: none;
      border-color: var(--primary);
      box-shadow: 0 0 0 3px rgba(0, 255, 65, 0.2);
    }
    
    input::placeholder {
      color: rgba(255, 255, 255, 0.3);
    }
    
    .captcha-container {
      display: flex;
      gap: 1rem;
      align-items: center;
    }
    
    .captcha-image {
      flex: 1;
      height: 60px;
      background: rgba(20, 20, 20, 0.7);
      border-radius: 8px;
      border: 1px solid rgba(255, 255, 255, 0.1);
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
    }
    
    .captcha-image img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    
    .captcha-refresh {
      background: rgba(20, 20, 20, 0.7);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 8px;
      color: var(--text);
      width: 50px;
      height: 50px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.3s ease;
    }
    
    .captcha-refresh:hover {
      background: rgba(0, 255, 65, 0.1);
      border-color: var(--primary);
      color: var(--primary);
    }
    
    .advanced-toggle {
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 1.5rem 0;
    }
    
    .advanced-toggle button {
      background: none;
      border: none;
      color: var(--primary);
      font-size: 0.95rem;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      transition: all 0.3s ease;
    }
    
    .advanced-toggle button:hover {
      background: rgba(0, 255, 65, 0.1);
    }
    
    .advanced-content {
      overflow: hidden;
      transition: max-height 0.5s ease, opacity 0.3s ease;
      max-height: 0;
      opacity: 0;
    }
    
    .advanced-content.show {
      max-height: 400px;
      opacity: 1;
    }
    
    .options-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1.5rem;
      margin-top: 1rem;
    }
    
    .option {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .checkbox-group {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }
    
    .checkbox-group input[type="checkbox"] {
      width: 20px;
      height: 20px;
      accent-color: var(--primary);
    }
    
    .checkbox-group label {
      margin: 0;
      color: var(--text);
      font-size: 0.95rem;
    }
    
    .btn {
      background: linear-gradient(135deg, var(--primary), var(--primary-dark));
      color: black;
      border: none;
      border-radius: 8px;
      padding: 1rem 2rem;
      font-size: 1.1rem;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.3s ease;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.75rem;
      width: 100%;
    }
    
    .btn:hover {
      transform: translateY(-2px);
      box-shadow: 0 5px 15px rgba(0, 255, 65, 0.3);
    }
    
    .btn:disabled {
      opacity: 0.7;
      cursor: not-allowed;
      transform: none;
      box-shadow: none;
    }
    
    #output {
      margin-top: 1.5rem;
      background: rgba(20, 20, 20, 0.7);
      border-radius: 8px;
      border: 1px solid rgba(0, 255, 65, 0.2);
      overflow: hidden;
    }
    
    .output-group {
      display: flex;
    }
    
    .output-group input {
      flex: 1;
      border: none;
      background: transparent;
      padding: 1rem;
      font-size: 1rem;
    }
    
    .copy-btn {
      background: rgba(0, 255, 65, 0.1);
      color: var(--primary);
      border: none;
      padding: 0 1.5rem;
      cursor: pointer;
      transition: all 0.3s ease;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    
    .copy-btn:hover {
      background: rgba(0, 255, 65, 0.2);
    }
    
    #error {
      display: flex;
      align-items: center;
      gap: 10px;
      background: rgba(255, 56, 96, 0.1);
      border: 1px solid var(--error);
      border-radius: 8px;
      padding: 1rem;
      margin-top: 1.5rem;
      color: var(--error);
      font-size: 0.95rem;
    }
    
    .features {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 1.5rem;
      margin-top: 1rem;
    }
    
    .feature {
      background: var(--bg-card);
      border-radius: 8px;
      padding: 1.5rem;
      border: 1px solid rgba(0, 255, 65, 0.1);
      transition: all 0.3s ease;
    }
    
    .feature:hover {
      border-color: var(--primary);
      transform: translateY(-5px);
    }
    
    .feature i {
      font-size: 2rem;
      color: var(--primary);
      margin-bottom: 1rem;
    }
    
    .feature h3 {
      font-size: 1.2rem;
      margin-bottom: 0.75rem;
      color: var(--text);
    }
    
    .feature p {
      color: var(--text-secondary);
      font-size: 0.95rem;
      line-height: 1.6;
    }
    
    footer {
      margin-top: 3rem;
      text-align: center;
      color: var(--text-secondary);
      font-size: 0.9rem;
      padding: 1rem;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
      width: 100%;
    }
    
    .matrix-bg {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      z-index: -1;
      opacity: 0.05;
      pointer-events: none;
      background: 
        linear-gradient(rgba(0, 255, 65, 0.05) 1px, transparent 1px),
        linear-gradient(90deg, rgba(0, 255, 65, 0.05) 1px, transparent 1px);
      background-size: 20px 20px;
    }
    
    @media (max-width: 768px) {
      .options-grid {
        grid-template-columns: 1fr;
      }
      
      .stats-container {
        grid-template-columns: 1fr;
      }
      
      .captcha-container {
        flex-direction: column;
      }
      
      .captcha-refresh {
        width: 100%;
        height: 50px;
      }
    }
  </style>
</head>
<body>
  <div class="matrix-bg"></div>
  
  <div class="container">
    <header>
      <div class="logo">
        <i class="fas fa-link logo-icon"></i>
        <h1>ktrk.ee</h1>
      </div>
      <p class="tagline">Premium URL Shortening Service</p>
      
      <div class="stats-container" id="stats">
        <div class="stat">
          <div class="stat-value">-</div>
          <div class="stat-label">Total Links</div>
        </div>
        <div class="stat">
          <div class="stat-value">-</div>
          <div class="stat-label">Active</div>
        </div>
        <div class="stat">
          <div class="stat-value">-</div>
          <div class="stat-label">Today</div>
        </div>
      </div>
    </header>
    
    <main class="card">
      <h2 class="card-title"><i class="fas fa-compress-alt"></i> Shorten Your URL</h2>
      
      <form id="shorten-form">
        <div class="form-group">
          <label for="original-url">Original URL</label>
          <div class="input-wrapper">
            <i class="fas fa-link"></i>
            <input 
              type="url" 
              id="original-url" 
              placeholder="https://example.com/long-url-path" 
              required
            >
          </div>
        </div>
        
        <div class="form-group">
          <label for="captcha-input">Verify You're Human</label>
          <div class="captcha-container">
            <div class="captcha-image" id="captcha-image">
              <div style="display: flex; align-items: center; justify-content: center; color: var(--text-secondary); font-size: 0.9rem;">
                Loading CAPTCHA...
              </div>
            </div>
            <div class="captcha-refresh" id="refresh-captcha">
              <i class="fas fa-sync-alt"></i>
            </div>
          </div>
          <div class="input-wrapper" style="margin-top: 0.75rem;">
            <i class="fas fa-shield-alt"></i>
            <input 
              type="text" 
              id="captcha-input" 
              placeholder="Enter the code above" 
              required
            >
          </div>
        </div>
        
        <div class="advanced-toggle">
          <button type="button" id="show-advanced">
            <i class="fas fa-chevron-down"></i> Show Advanced Options
          </button>
        </div>
        
        <div class="advanced-content" id="advanced-content">
          <div class="options-grid">
            <div class="option">
              <label>Expiration Settings</label>
              <div class="checkbox-group">
                <input type="checkbox" id="permanent-checkbox">
                <label for="permanent-checkbox">Permanent Link</label>
              </div>
            </div>
            
            <div class="option">
              <label for="ttl-input">Link Lifetime (Days)</label>
              <input 
                type="number" 
                id="ttl-input" 
                min="1" 
                max="365" 
                value="30"
              >
            </div>
          </div>
          
          <div class="form-group" style="margin-top: 1.5rem;">
            <label for="metadata-input">Custom Metadata (Optional)</label>
            <div class="input-wrapper">
              <i class="fas fa-tag"></i>
              <input 
                type="text" 
                id="metadata-input" 
                placeholder="campaign, source, etc."
              >
            </div>
          </div>
        </div>
        
        <button type="submit" class="btn" id="shorten-btn">
          <i class="fas fa-link"></i> Shorten URL
        </button>
      </form>
      
      <div id="output" style="display: none;">
        <div class="output-group">
          <input type="text" id="shortened-url" readonly>
          <button class="copy-btn" id="copy-btn">
            <i class="fas fa-copy"></i> Copy
          </button>
        </div>
      </div>
      
      <div id="error" style="display: none;"></div>
    </main>
    
  </div>
  
  <footer>
    <p>© 2023 ktrk.ee | Premium URL Shortening Service | v2.1.0</p>
  </footer>

  <script>
    document.addEventListener('DOMContentLoaded', () => {
      // API base URL - измените на ваш реальный адрес сервера
      const API_BASE_URL = 'http://localhost:8080/api';

      // State
      let state = {
        originalUrl: '',
        captchaId: '',
        captchaSolution: '',
        shortenedUrl: '',
        loading: false,
        error: '',
        isPermanent: false,
        ttlDays: 30,
        metadata: '',
        showAdvanced: false,
        stats: {
          totalShortened: 0,
          activeLinks: 0,
          todayLinks: 0
        }
      };

      // Elements
      const elements = {
        originalUrl: document.getElementById('original-url'),
        captchaInput: document.getElementById('captcha-input'),
        captchaImage: document.getElementById('captcha-image'),
        permanentCheckbox: document.getElementById('permanent-checkbox'),
        ttlInput: document.getElementById('ttl-input'),
        metadataInput: document.getElementById('metadata-input'),
        showAdvancedBtn: document.getElementById('show-advanced'),
        advancedContent: document.getElementById('advanced-content'),
        shortenBtn: document.getElementById('shorten-btn'),
        form: document.getElementById('shorten-form'),
        output: document.getElementById('output'),
        shortenedUrl: document.getElementById('shortened-url'),
        copyBtn: document.getElementById('copy-btn'),
        error: document.getElementById('error'),
        refreshCaptcha: document.getElementById('refresh-captcha'),
        stats: document.getElementById('stats')
      };

      // API calls
      async function apiCall(endpoint, options = {}) {
        try {
          const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            headers: {
              'Content-Type': 'application/json',
              ...options.headers
            },
            ...options
          });

          if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'API request failed');
          }

          return await response.json();
        } catch (error) {
          console.error('API Error:', error);
          throw error;
        }
      }

      // Load CAPTCHA
      async function loadCaptcha() {
        try {
          const data = await apiCall('/captcha');
          state.captchaId = data.captcha_id;
          
          // Load CAPTCHA image
          const captchaImg = new Image();
          captchaImg.onload = () => {
            elements.captchaImage.innerHTML = '';
            elements.captchaImage.appendChild(captchaImg);
          };
          captchaImg.onerror = () => {
            elements.captchaImage.innerHTML = `
              <div style="color: var(--error); text-align: center;">
                Failed to load CAPTCHA
              </div>
            `;
          };
          captchaImg.src = `/captcha/${state.captchaId}`;
          captchaImg.style.width = '100%';
          captchaImg.style.height = '100%';
          captchaImg.style.objectFit = 'cover';
          
        } catch (error) {
          console.error('Failed to load CAPTCHA:', error);
          elements.captchaImage.innerHTML = `
            <div style="color: var(--error); text-align: center;">
              Failed to load CAPTCHA
            </div>
          `;
        }
      }

      // Refresh captcha
      async function refreshCaptcha() {
        const refreshIcon = elements.refreshCaptcha.querySelector('i');
        refreshIcon.classList.add('fa-spin');
        
        try {
          await loadCaptcha();
          elements.captchaInput.value = '';
          state.captchaSolution = '';
        } finally {
          refreshIcon.classList.remove('fa-spin');
        }
      }

      // Load stats
      async function loadStats() {
        try {
          const stats = await apiCall('/stats');
          state.stats = stats;
          updateStats();
        } catch (error) {
          console.error('Failed to load stats:', error);
          // Keep default values or show error
        }
      }

      // Shorten URL
      async function shortenUrl(e) {
        e.preventDefault();
        
        if (!state.originalUrl) {
          showError('Please enter a valid URL');
          return;
        }
        
        if (!state.captchaSolution) {
          showError('Please complete the CAPTCHA');
          return;
        }
        
        state.loading = true;
        state.error = '';
        updateUI();
        
        try {
          const requestData = {
            originalUrl: state.originalUrl,
            captchaId: state.captchaId,
            captchaToken: state.captchaSolution,
            metadata: state.metadata || '',
            isPermanent: state.isPermanent
          };

          // Add TTL only if not permanent
          if (!state.isPermanent) {
            requestData.ttl = state.ttlDays * 24 * 60 * 60; // Convert days to seconds
          }

          const result = await apiCall('/shorten', {
            method: 'POST',
            body: JSON.stringify(requestData)
          });
          
          state.shortenedUrl = result.shortUrl;
          elements.output.style.display = 'block';
          elements.shortenedUrl.value = state.shortenedUrl;
          
          // Reset form
          elements.originalUrl.value = '';
          elements.captchaInput.value = '';
          state.originalUrl = '';
          state.captchaSolution = '';
          
          // Reload captcha and stats
          await loadCaptcha();
          await loadStats();
          
        } catch (error) {
          showError(error.message || 'Failed to shorten URL. Please try again.');
        } finally {
          state.loading = false;
          updateUI();
        }
      }

      // Update stats display
      function updateStats() {
        const statValues = elements.stats.querySelectorAll('.stat-value');
        if (statValues.length === 3) {
          statValues[0].textContent = state.stats.totalShortened.toLocaleString();
          statValues[1].textContent = state.stats.activeLinks.toLocaleString();
          statValues[2].textContent = state.stats.todayLinks.toLocaleString();
        }
      }
      
      // Show error
      function showError(message) {
        state.error = message;
        elements.error.innerHTML = `<i class="fas fa-exclamation-triangle"></i> ${message}`;
        elements.error.style.display = 'flex';
        setTimeout(() => {
          elements.error.style.display = 'none';
        }, 5000);
      }
      
      // Copy to clipboard
      function copyToClipboard() {
        navigator.clipboard.writeText(state.shortenedUrl).then(() => {
          elements.copyBtn.innerHTML = '<i class="fas fa-check"></i> Copied!';
          setTimeout(() => {
            elements.copyBtn.innerHTML = '<i class="fas fa-copy"></i> Copy';
          }, 2000);
        }).catch(() => {
          // Fallback for older browsers
          elements.shortenedUrl.select();
          document.execCommand('copy');
          elements.copyBtn.innerHTML = '<i class="fas fa-check"></i> Copied!';
          setTimeout(() => {
            elements.copyBtn.innerHTML = '<i class="fas fa-copy"></i> Copy';
          }, 2000);
        });
      }
      
      // Update UI based on state
      function updateUI() {
        elements.shortenBtn.disabled = state.loading;
        elements.shortenBtn.innerHTML = state.loading 
          ? '<i class="fas fa-spinner fa-spin"></i> Shortening...' 
          : '<i class="fas fa-link"></i> Shorten URL';
        
        elements.ttlInput.disabled = state.isPermanent;
        elements.ttlInput.style.opacity = state.isPermanent ? '0.5' : '1';
      }
      
      // Toggle advanced options
      function toggleAdvanced() {
        state.showAdvanced = !state.showAdvanced;
        
        if (state.showAdvanced) {
          elements.advancedContent.classList.add('show');
          elements.showAdvancedBtn.innerHTML = '<i class="fas fa-chevron-up"></i> Hide Advanced Options';
        } else {
          elements.advancedContent.classList.remove('show');
          elements.showAdvancedBtn.innerHTML = '<i class="fas fa-chevron-down"></i> Show Advanced Options';
        }
      }
      
      // Event listeners
      elements.form.addEventListener('submit', shortenUrl);
      elements.refreshCaptcha.addEventListener('click', refreshCaptcha);
      elements.permanentCheckbox.addEventListener('change', () => {
        state.isPermanent = elements.permanentCheckbox.checked;
        updateUI();
      });
      elements.copyBtn.addEventListener('click', copyToClipboard);
      elements.showAdvancedBtn.addEventListener('click', toggleAdvanced);
      
      // Input bindings
      elements.originalUrl.addEventListener('input', (e) => state.originalUrl = e.target.value);
      elements.captchaInput.addEventListener('input', (e) => state.captchaSolution = e.target.value);
      elements.ttlInput.addEventListener('input', (e) => state.ttlDays = parseInt(e.target.value) || 30);
      elements.metadataInput.addEventListener('input', (e) => state.metadata = e.target.value);
      
      // Initialize
      loadCaptcha();
      loadStats();
    });
  </script>
</body>
</html>
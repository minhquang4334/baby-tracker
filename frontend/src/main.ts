import '../css/main.css';
import '../css/components.css';
import '../css/onboarding.css';
import '../css/dashboard.css';
import '../css/modal.css';
import '../css/history.css';
import '../css/growth-chart.css';
import '../css/analytics.css';

import { addRoute, initRouter, navigate } from './router';
import { api } from './api';
import { state } from './state';
import { renderOnboarding } from './components/onboarding';
import { renderDashboard } from './components/dashboard';
import { renderHistory } from './components/history';
import { renderGrowthScreen } from './components/growth-chart';
import { renderGuideScreen } from './components/guide';
import { renderAnalyticsScreen } from './components/analytics';

// Register routes
addRoute('/onboarding', renderOnboarding);
addRoute('/dashboard', renderDashboard);
addRoute('/history', renderHistory);
addRoute('/analytics', renderAnalyticsScreen);
addRoute('/growth', renderGrowthScreen);
addRoute('/guide', renderGuideScreen);

async function bootstrap() {
  const app = document.getElementById('app');
  if (!app) return;

  try {
    const child = await api.getChild();
    if (child && 'id' in child) {
      state.child.set(child);
      // If no hash or hash is onboarding, go to dashboard
      const hash = window.location.hash.slice(1);
      if (!hash || hash === '/onboarding') {
        navigate('/dashboard');
      }
    } else {
      navigate('/onboarding');
    }
  } catch {
    navigate('/onboarding');
  }

  initRouter(app);
}

bootstrap();

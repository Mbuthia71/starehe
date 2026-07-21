import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.siohioma.banking',
  appName: 'Kindred Mobile',
  webDir: '.output/public',
  server: {
    url: 'https://banking.siohioma.com',
    cleartext: false,
    allowNavigation: ['*']
  }
};

export default config;

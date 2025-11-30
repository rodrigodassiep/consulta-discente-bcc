/// <reference types="vite-plugin-pwa/info" />
/// <reference types="vite-plugin-pwa/pwa-assets" />

declare module 'virtual:pwa-info' {
	export const pwaInfo: {
		webManifest: {
			href: string;
			useCredentials: boolean;
		};
	} | undefined;
}

declare module 'virtual:pwa-register' {
	export function registerSW(options?: {
		immediate?: boolean;
		onNeedRefresh?: () => void;
		onOfflineReady?: () => void;
		onRegistered?: (registration: ServiceWorkerRegistration | undefined) => void;
		onRegisterError?: (error: Error) => void;
	}): (reloadPage?: boolean) => Promise<void>;
}

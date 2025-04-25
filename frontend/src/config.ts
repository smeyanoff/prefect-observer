let host: string;

if (process.env.VUE_APP_BACKEND_API_HOST) {
  host = process.env.VUE_APP_BACKEND_API_HOST;
} else {
  host = "http://localhost:8180";
}

const config = {
  backendApiHost: host,
};

export default config;

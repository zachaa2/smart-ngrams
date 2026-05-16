import { resolve } from "path"
import { defineConfig } from "vite"

export default defineConfig({
    base: "/smart-ngrams/",
    build: {
        rolldownOptions:  {
            input: {
                main: resolve(__dirname, "index.html"),
                algorithm: resolve(__dirname, "algorithm.html"),
            }
        }
    }
})
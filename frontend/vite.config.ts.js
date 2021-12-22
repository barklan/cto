// vite.config.ts
import path from "path";
import { defineConfig } from "vite";
import Vue from "@vitejs/plugin-vue";
import Pages from "vite-plugin-pages";
import Layouts from "vite-plugin-vue-layouts";
import Icons from "unplugin-icons/vite";
import IconsResolver from "unplugin-icons/resolver";
import Components from "unplugin-vue-components/vite";
import AutoImport from "unplugin-auto-import/vite";
import Markdown from "vite-plugin-md";
import WindiCSS from "vite-plugin-windicss";
import { VitePWA } from "vite-plugin-pwa";
import VueI18n from "@intlify/vite-plugin-vue-i18n";
import Inspect from "vite-plugin-inspect";
import Prism from "markdown-it-prism";
import LinkAttributes from "markdown-it-link-attributes";
var markdownWrapperClasses = "prose prose-sm m-auto text-left";
var vite_config_default = defineConfig({
  build: {
    rollupOptions: {
      external: [
        "workbox-window"
      ]
    }
  },
  resolve: {
    alias: {
      "~/": `${path.resolve("/home/barklan/dev/cto/frontend", "src")}/`
    }
  },
  plugins: [
    Vue({
      include: [/\.vue$/, /\.md$/]
    }),
    Pages({
      extensions: ["vue", "md"]
    }),
    Layouts(),
    AutoImport({
      imports: [
        "vue",
        "vue-router",
        "vue-i18n",
        "@vueuse/head",
        "@vueuse/core"
      ],
      dts: "src/auto-imports.d.ts"
    }),
    Components({
      extensions: ["vue", "md"],
      include: [/\.vue$/, /\.vue\?vue/, /\.md$/],
      resolvers: [
        IconsResolver({
          componentPrefix: ""
        })
      ],
      dts: "src/components.d.ts"
    }),
    Icons({
      autoInstall: true
    }),
    WindiCSS({
      safelist: markdownWrapperClasses
    }),
    Markdown({
      wrapperClasses: markdownWrapperClasses,
      headEnabled: true,
      markdownItSetup(md) {
        md.use(Prism);
        md.use(LinkAttributes, {
          pattern: /^https?:\/\//,
          attrs: {
            target: "_blank",
            rel: "noopener"
          }
        });
      }
    }),
    VitePWA({
      registerType: "autoUpdate",
      includeAssets: ["favicon.svg", "robots.txt", "safari-pinned-tab.svg"],
      manifest: {
        name: "CTO",
        short_name: "CTO",
        theme_color: "#ffffff",
        icons: [
          {
            src: "/pwa-192x192.png",
            sizes: "192x192",
            type: "image/png"
          },
          {
            src: "/pwa-512x512.png",
            sizes: "512x512",
            type: "image/png"
          },
          {
            src: "/pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
            purpose: "any maskable"
          }
        ]
      }
    }),
    VueI18n({
      runtimeOnly: true,
      compositionOnly: true,
      include: [path.resolve("/home/barklan/dev/cto/frontend", "locales/**")]
    }),
    Inspect({
      enabled: false
    })
  ],
  server: {
    fs: {
      strict: true
    }
  },
  ssgOptions: {
    script: "async",
    formatting: "minify"
  },
  optimizeDeps: {
    include: [
      "vue",
      "vue-router",
      "@vueuse/core",
      "@vueuse/head"
    ],
    exclude: [
      "vue-demi"
    ]
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImltcG9ydCBwYXRoIGZyb20gJ3BhdGgnXG5pbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tICd2aXRlJ1xuaW1wb3J0IFZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXG5pbXBvcnQgUGFnZXMgZnJvbSAndml0ZS1wbHVnaW4tcGFnZXMnXG5pbXBvcnQgTGF5b3V0cyBmcm9tICd2aXRlLXBsdWdpbi12dWUtbGF5b3V0cydcbmltcG9ydCBJY29ucyBmcm9tICd1bnBsdWdpbi1pY29ucy92aXRlJ1xuaW1wb3J0IEljb25zUmVzb2x2ZXIgZnJvbSAndW5wbHVnaW4taWNvbnMvcmVzb2x2ZXInXG5pbXBvcnQgQ29tcG9uZW50cyBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy92aXRlJ1xuaW1wb3J0IEF1dG9JbXBvcnQgZnJvbSAndW5wbHVnaW4tYXV0by1pbXBvcnQvdml0ZSdcbmltcG9ydCBNYXJrZG93biBmcm9tICd2aXRlLXBsdWdpbi1tZCdcbmltcG9ydCBXaW5kaUNTUyBmcm9tICd2aXRlLXBsdWdpbi13aW5kaWNzcydcbmltcG9ydCB7IFZpdGVQV0EgfSBmcm9tICd2aXRlLXBsdWdpbi1wd2EnXG5pbXBvcnQgVnVlSTE4biBmcm9tICdAaW50bGlmeS92aXRlLXBsdWdpbi12dWUtaTE4bidcbmltcG9ydCBJbnNwZWN0IGZyb20gJ3ZpdGUtcGx1Z2luLWluc3BlY3QnXG5pbXBvcnQgUHJpc20gZnJvbSAnbWFya2Rvd24taXQtcHJpc20nXG5pbXBvcnQgTGlua0F0dHJpYnV0ZXMgZnJvbSAnbWFya2Rvd24taXQtbGluay1hdHRyaWJ1dGVzJ1xuXG5jb25zdCBtYXJrZG93bldyYXBwZXJDbGFzc2VzID0gJ3Byb3NlIHByb3NlLXNtIG0tYXV0byB0ZXh0LWxlZnQnXG5cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XG4gIGJ1aWxkOiB7XG4gICAgcm9sbHVwT3B0aW9uczoge1xuICAgICAgZXh0ZXJuYWw6IFtcbiAgICAgICAgXCJ3b3JrYm94LXdpbmRvd1wiLFxuICAgICAgXSxcbiAgICB9LFxuICB9LFxuXG4gIHJlc29sdmU6IHtcbiAgICBhbGlhczoge1xuICAgICAgJ34vJzogYCR7cGF0aC5yZXNvbHZlKFwiL2hvbWUvYmFya2xhbi9kZXYvY3RvL2Zyb250ZW5kXCIsICdzcmMnKX0vYCxcbiAgICB9LFxuICB9LFxuICBwbHVnaW5zOiBbXG4gICAgVnVlKHtcbiAgICAgIGluY2x1ZGU6IFsvXFwudnVlJC8sIC9cXC5tZCQvXSxcbiAgICB9KSxcblxuICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9oYW5ub2VydS92aXRlLXBsdWdpbi1wYWdlc1xuICAgIFBhZ2VzKHtcbiAgICAgIGV4dGVuc2lvbnM6IFsndnVlJywgJ21kJ10sXG4gICAgfSksXG5cbiAgICAvLyBodHRwczovL2dpdGh1Yi5jb20vSm9obkNhbXBpb25Kci92aXRlLXBsdWdpbi12dWUtbGF5b3V0c1xuICAgIExheW91dHMoKSxcblxuICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9hbnRmdS91bnBsdWdpbi1hdXRvLWltcG9ydFxuICAgIEF1dG9JbXBvcnQoe1xuICAgICAgaW1wb3J0czogW1xuICAgICAgICAndnVlJyxcbiAgICAgICAgJ3Z1ZS1yb3V0ZXInLFxuICAgICAgICAndnVlLWkxOG4nLFxuICAgICAgICAnQHZ1ZXVzZS9oZWFkJyxcbiAgICAgICAgJ0B2dWV1c2UvY29yZScsXG4gICAgICBdLFxuICAgICAgZHRzOiAnc3JjL2F1dG8taW1wb3J0cy5kLnRzJyxcbiAgICB9KSxcblxuICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9hbnRmdS91bnBsdWdpbi12dWUtY29tcG9uZW50c1xuICAgIENvbXBvbmVudHMoe1xuICAgICAgLy8gYWxsb3cgYXV0byBsb2FkIG1hcmtkb3duIGNvbXBvbmVudHMgdW5kZXIgYC4vc3JjL2NvbXBvbmVudHMvYFxuICAgICAgZXh0ZW5zaW9uczogWyd2dWUnLCAnbWQnXSxcblxuICAgICAgLy8gYWxsb3cgYXV0byBpbXBvcnQgYW5kIHJlZ2lzdGVyIGNvbXBvbmVudHMgdXNlZCBpbiBtYXJrZG93blxuICAgICAgaW5jbHVkZTogWy9cXC52dWUkLywgL1xcLnZ1ZVxcP3Z1ZS8sIC9cXC5tZCQvXSxcblxuICAgICAgLy8gY3VzdG9tIHJlc29sdmVyc1xuICAgICAgcmVzb2x2ZXJzOiBbXG4gICAgICAgIC8vIGF1dG8gaW1wb3J0IGljb25zXG4gICAgICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9hbnRmdS91bnBsdWdpbi1pY29uc1xuICAgICAgICBJY29uc1Jlc29sdmVyKHtcbiAgICAgICAgICBjb21wb25lbnRQcmVmaXg6ICcnLFxuICAgICAgICAgIC8vIGVuYWJsZWRDb2xsZWN0aW9uczogWydjYXJib24nXVxuICAgICAgICB9KSxcbiAgICAgIF0sXG5cbiAgICAgIGR0czogJ3NyYy9jb21wb25lbnRzLmQudHMnLFxuICAgIH0pLFxuXG4gICAgLy8gaHR0cHM6Ly9naXRodWIuY29tL2FudGZ1L3VucGx1Z2luLWljb25zXG4gICAgSWNvbnMoe1xuICAgICAgYXV0b0luc3RhbGw6IHRydWUsXG4gICAgfSksXG5cbiAgICAvLyBodHRwczovL2dpdGh1Yi5jb20vYW50ZnUvdml0ZS1wbHVnaW4td2luZGljc3NcbiAgICBXaW5kaUNTUyh7XG4gICAgICBzYWZlbGlzdDogbWFya2Rvd25XcmFwcGVyQ2xhc3NlcyxcbiAgICB9KSxcblxuICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9hbnRmdS92aXRlLXBsdWdpbi1tZFxuICAgIC8vIERvbid0IG5lZWQgdGhpcz8gVHJ5IHZpdGVzc2UtbGl0ZTogaHR0cHM6Ly9naXRodWIuY29tL2FudGZ1L3ZpdGVzc2UtbGl0ZVxuICAgIE1hcmtkb3duKHtcbiAgICAgIHdyYXBwZXJDbGFzc2VzOiBtYXJrZG93bldyYXBwZXJDbGFzc2VzLFxuICAgICAgaGVhZEVuYWJsZWQ6IHRydWUsXG4gICAgICBtYXJrZG93bkl0U2V0dXAobWQpIHtcbiAgICAgICAgLy8gaHR0cHM6Ly9wcmlzbWpzLmNvbS9cbiAgICAgICAgLy8gQHRzLWV4cGVjdC1lcnJvciB0eXBlcyBtaXNtYXRjaFxuICAgICAgICBtZC51c2UoUHJpc20pXG4gICAgICAgIC8vIEB0cy1leHBlY3QtZXJyb3IgdHlwZXMgbWlzbWF0Y2hcbiAgICAgICAgbWQudXNlKExpbmtBdHRyaWJ1dGVzLCB7XG4gICAgICAgICAgcGF0dGVybjogL15odHRwcz86XFwvXFwvLyxcbiAgICAgICAgICBhdHRyczoge1xuICAgICAgICAgICAgdGFyZ2V0OiAnX2JsYW5rJyxcbiAgICAgICAgICAgIHJlbDogJ25vb3BlbmVyJyxcbiAgICAgICAgICB9LFxuICAgICAgICB9KVxuICAgICAgfSxcbiAgICB9KSxcblxuICAgIC8vIGh0dHBzOi8vZ2l0aHViLmNvbS9hbnRmdS92aXRlLXBsdWdpbi1wd2FcbiAgICBWaXRlUFdBKHtcbiAgICAgIHJlZ2lzdGVyVHlwZTogJ2F1dG9VcGRhdGUnLFxuICAgICAgaW5jbHVkZUFzc2V0czogWydmYXZpY29uLnN2ZycsICdyb2JvdHMudHh0JywgJ3NhZmFyaS1waW5uZWQtdGFiLnN2ZyddLFxuICAgICAgbWFuaWZlc3Q6IHtcbiAgICAgICAgbmFtZTogJ0NUTycsXG4gICAgICAgIHNob3J0X25hbWU6ICdDVE8nLFxuICAgICAgICB0aGVtZV9jb2xvcjogJyNmZmZmZmYnLFxuICAgICAgICBpY29uczogW1xuICAgICAgICAgIHtcbiAgICAgICAgICAgIHNyYzogJy9wd2EtMTkyeDE5Mi5wbmcnLFxuICAgICAgICAgICAgc2l6ZXM6ICcxOTJ4MTkyJyxcbiAgICAgICAgICAgIHR5cGU6ICdpbWFnZS9wbmcnLFxuICAgICAgICAgIH0sXG4gICAgICAgICAge1xuICAgICAgICAgICAgc3JjOiAnL3B3YS01MTJ4NTEyLnBuZycsXG4gICAgICAgICAgICBzaXplczogJzUxMng1MTInLFxuICAgICAgICAgICAgdHlwZTogJ2ltYWdlL3BuZycsXG4gICAgICAgICAgfSxcbiAgICAgICAgICB7XG4gICAgICAgICAgICBzcmM6ICcvcHdhLTUxMng1MTIucG5nJyxcbiAgICAgICAgICAgIHNpemVzOiAnNTEyeDUxMicsXG4gICAgICAgICAgICB0eXBlOiAnaW1hZ2UvcG5nJyxcbiAgICAgICAgICAgIHB1cnBvc2U6ICdhbnkgbWFza2FibGUnLFxuICAgICAgICAgIH0sXG4gICAgICAgIF0sXG4gICAgICB9LFxuICAgIH0pLFxuXG4gICAgLy8gaHR0cHM6Ly9naXRodWIuY29tL2ludGxpZnkvYnVuZGxlLXRvb2xzL3RyZWUvbWFpbi9wYWNrYWdlcy92aXRlLXBsdWdpbi12dWUtaTE4blxuICAgIFZ1ZUkxOG4oe1xuICAgICAgcnVudGltZU9ubHk6IHRydWUsXG4gICAgICBjb21wb3NpdGlvbk9ubHk6IHRydWUsXG4gICAgICBpbmNsdWRlOiBbcGF0aC5yZXNvbHZlKFwiL2hvbWUvYmFya2xhbi9kZXYvY3RvL2Zyb250ZW5kXCIsICdsb2NhbGVzLyoqJyldLFxuICAgIH0pLFxuXG4gICAgLy8gaHR0cHM6Ly9naXRodWIuY29tL2FudGZ1L3ZpdGUtcGx1Z2luLWluc3BlY3RcbiAgICBJbnNwZWN0KHtcbiAgICAgIC8vIGNoYW5nZSB0aGlzIHRvIGVuYWJsZSBpbnNwZWN0IGZvciBkZWJ1Z2dpbmdcbiAgICAgIGVuYWJsZWQ6IGZhbHNlLFxuICAgIH0pLFxuICBdLFxuXG4gIHNlcnZlcjoge1xuICAgIGZzOiB7XG4gICAgICBzdHJpY3Q6IHRydWUsXG4gICAgfSxcbiAgfSxcblxuICAvLyBodHRwczovL2dpdGh1Yi5jb20vYW50ZnUvdml0ZS1zc2dcbiAgc3NnT3B0aW9uczoge1xuICAgIHNjcmlwdDogJ2FzeW5jJyxcbiAgICBmb3JtYXR0aW5nOiAnbWluaWZ5JyxcbiAgfSxcblxuICBvcHRpbWl6ZURlcHM6IHtcbiAgICBpbmNsdWRlOiBbXG4gICAgICAndnVlJyxcbiAgICAgICd2dWUtcm91dGVyJyxcbiAgICAgICdAdnVldXNlL2NvcmUnLFxuICAgICAgJ0B2dWV1c2UvaGVhZCcsXG4gICAgXSxcbiAgICBleGNsdWRlOiBbXG4gICAgICAndnVlLWRlbWknLFxuICAgIF0sXG4gIH0sXG59KVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUFBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBRUEsSUFBTSx5QkFBeUI7QUFFL0IsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsT0FBTztBQUFBLElBQ0wsZUFBZTtBQUFBLE1BQ2IsVUFBVTtBQUFBLFFBQ1I7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtOLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLE1BQU0sR0FBRyxLQUFLLFFBQVEsa0NBQWtDO0FBQUE7QUFBQTtBQUFBLEVBRzVELFNBQVM7QUFBQSxJQUNQLElBQUk7QUFBQSxNQUNGLFNBQVMsQ0FBQyxVQUFVO0FBQUE7QUFBQSxJQUl0QixNQUFNO0FBQUEsTUFDSixZQUFZLENBQUMsT0FBTztBQUFBO0FBQUEsSUFJdEI7QUFBQSxJQUdBLFdBQVc7QUFBQSxNQUNULFNBQVM7QUFBQSxRQUNQO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBO0FBQUEsTUFFRixLQUFLO0FBQUE7QUFBQSxJQUlQLFdBQVc7QUFBQSxNQUVULFlBQVksQ0FBQyxPQUFPO0FBQUEsTUFHcEIsU0FBUyxDQUFDLFVBQVUsY0FBYztBQUFBLE1BR2xDLFdBQVc7QUFBQSxRQUdULGNBQWM7QUFBQSxVQUNaLGlCQUFpQjtBQUFBO0FBQUE7QUFBQSxNQUtyQixLQUFLO0FBQUE7QUFBQSxJQUlQLE1BQU07QUFBQSxNQUNKLGFBQWE7QUFBQTtBQUFBLElBSWYsU0FBUztBQUFBLE1BQ1AsVUFBVTtBQUFBO0FBQUEsSUFLWixTQUFTO0FBQUEsTUFDUCxnQkFBZ0I7QUFBQSxNQUNoQixhQUFhO0FBQUEsTUFDYixnQkFBZ0IsSUFBSTtBQUdsQixXQUFHLElBQUk7QUFFUCxXQUFHLElBQUksZ0JBQWdCO0FBQUEsVUFDckIsU0FBUztBQUFBLFVBQ1QsT0FBTztBQUFBLFlBQ0wsUUFBUTtBQUFBLFlBQ1IsS0FBSztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsSUFPYixRQUFRO0FBQUEsTUFDTixjQUFjO0FBQUEsTUFDZCxlQUFlLENBQUMsZUFBZSxjQUFjO0FBQUEsTUFDN0MsVUFBVTtBQUFBLFFBQ1IsTUFBTTtBQUFBLFFBQ04sWUFBWTtBQUFBLFFBQ1osYUFBYTtBQUFBLFFBQ2IsT0FBTztBQUFBLFVBQ0w7QUFBQSxZQUNFLEtBQUs7QUFBQSxZQUNMLE9BQU87QUFBQSxZQUNQLE1BQU07QUFBQTtBQUFBLFVBRVI7QUFBQSxZQUNFLEtBQUs7QUFBQSxZQUNMLE9BQU87QUFBQSxZQUNQLE1BQU07QUFBQTtBQUFBLFVBRVI7QUFBQSxZQUNFLEtBQUs7QUFBQSxZQUNMLE9BQU87QUFBQSxZQUNQLE1BQU07QUFBQSxZQUNOLFNBQVM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLElBT2pCLFFBQVE7QUFBQSxNQUNOLGFBQWE7QUFBQSxNQUNiLGlCQUFpQjtBQUFBLE1BQ2pCLFNBQVMsQ0FBQyxLQUFLLFFBQVEsa0NBQWtDO0FBQUE7QUFBQSxJQUkzRCxRQUFRO0FBQUEsTUFFTixTQUFTO0FBQUE7QUFBQTtBQUFBLEVBSWIsUUFBUTtBQUFBLElBQ04sSUFBSTtBQUFBLE1BQ0YsUUFBUTtBQUFBO0FBQUE7QUFBQSxFQUtaLFlBQVk7QUFBQSxJQUNWLFFBQVE7QUFBQSxJQUNSLFlBQVk7QUFBQTtBQUFBLEVBR2QsY0FBYztBQUFBLElBQ1osU0FBUztBQUFBLE1BQ1A7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQTtBQUFBLElBRUYsU0FBUztBQUFBLE1BQ1A7QUFBQTtBQUFBO0FBQUE7IiwKICAibmFtZXMiOiBbXQp9Cg==

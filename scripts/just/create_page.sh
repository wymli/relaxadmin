
cat  <<EOF 
{{define "page_${PAGE}_html"}}

<div class="h-full w-full">
    <div
        class="container font-sm mx-auto w-full p-2 block bg-gray-100 border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700">
        <label class="block text-lg font-medium text-gray-900 dark:text-white">Hello World</label>
        <label class="block text-sm text-gray-900"> hello world </label>
        <br />
        <div class="">
            Hello World
        </div>
        <br />
    </div>
    <div id="rsp"
        class="container mt-2 mx-auto h-auto w-full p-2 block bg-gray-100 border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700">
    </div>
</div>

{{end}}


{{define "page_${PAGE}_js"}}

<script>
    var app = {
        init() {
            \$(function () {})
        },
        getRegisterURL() {
            return "/page_$PAGE"
        },
        render() {
            global.setMainContent(\`{{ template "page_${PAGE}_html" . }}\`)
            \$(function () {
                app.init()
            })
        },
        getSidebar() {
            return \`<a class="flex tab items-center px-3 py-2 text-gray-600 transition-colors duration-300 transform rounded-lg dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 dark:hover:text-gray-200 hover:text-gray-700"
                href="#\${this.getRegisterURL()}">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                </svg>

                <span class="mx-2 text-sm font-medium capitalize">${PAGE}</span>
            </a>\`
        },
    }

    global.router.registerRouteApp(app.getRegisterURL(), app)
    // global.router.setDefaultRoute(app.getRegisterURL())
</script>

{{end}}

EOF
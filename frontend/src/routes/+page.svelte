<script>
    import Sidebar from "../components/sidebar.svelte";
</script>

<div class="flex flex-row">
    <div class="h-screen justify-center md:w-1/6 sm:w-2/6">
        <Sidebar/>
    </div>
    <div class="h-screen w-5/6">
        <div class="mockup-browser border border-base-300 h-screen">
            <div class="mockup-browser-toolbar">
              <div class="input border border-base-300">https://daisyui.com</div>
            </div>
            <div class="flex justify-center px-4 py-16 border-t border-base-300">Hello!</div>
        </div>
    </div>
</div>
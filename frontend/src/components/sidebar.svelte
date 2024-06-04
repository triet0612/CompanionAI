<script>
    async function createChat() {
        let new_chat = await fetch("http://localhost:8000/api/v1/chat", {
            method: "POST",
            credentials: "include"
        }).then(res => {
            if (res.status !== 200) {
                return {}
            }
            return res.json()
        }).catch(err => {
            console.log(err);
            alert(err)
            return {}
        })
        console.log(new_chat);
    }
</script>

<div class="h-full w-full bg-secondary pt-4 px-4 overflow-y-scroll">
    <div class="flex w-full">
        <button class="btn w-full mx-auto text-lg" on:click={async()=>{await createChat()}}>
            <span class="material-symbols-outlined">
                add
            </span>
        </button>
    </div>
    <div class="w-full">
        <ul class="w-full">
            <li class="flex flex-row w-full py-2 gap-2 justify-between">
                <button class="btn btn-primary w-4/6">Title</button>
                <button class="btn btn-primary w-1/6">
                    <span class="material-symbols-outlined">
                        delete
                    </span>
                </button>
            </li>
        </ul>
    </div>
</div>
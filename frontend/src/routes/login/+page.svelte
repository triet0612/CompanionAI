<script>
    let email = ""
    let password = ""
    async function login() {
        let res = await fetch("http://localhost:8000"+"/api/v1/login", {
            method: "POST",
            credentials:  "include",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        }).then(res => {
            if (res.status !== 200) {
                alert("Unauthorized")
            }
            return "ok"
        }).catch(err => {
            console.log(err);
            return "error"
        })
        if (res === "ok") {
            location.replace("/")
        }
    }
</script>

<div class="relative flex flex-col justify-center h-screen overflow-hidden">
    <div class="w-full p-6 m-auto bg-base-300 rounded-md shadow-md lg:max-w-lg">
        <h1 class="text-3xl font-semibold text-center text-neutral">Login</h1>
        <div>
            <label class="label text-lg text-primary-content label-text" for="email">Email</label>
            <input bind:value={email} name="email" type="text" placeholder="Email Address" class="w-full input input-bordered input-primary" />
        </div>
        <div>
            <label class="label text-lg text-primary-content label-text" for="password">Password</label>
            <input bind:value={password} name="password" type="password" placeholder="Enter Password" class="w-full input input-bordered input-primary" />
        </div>
        <div>
            <a href="/register" class="label label-text text-base text-primary-content hover:underline hover:text-blue-600">Sign in?</a>
        </div>
        <div class="w-full flex flex-row">
            <button on:click={async ()=>{await login()}} class="btn btn-primary w-1/2 mx-auto justify-center">Login</button>
        </div>
    </div>
</div>

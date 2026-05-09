const go = new Go();

const result = await WebAssembly.instantiateStreaming(
  fetch("/main.wasm"),
  go.importObject,
);

go.run(result.instance);

document.querySelector("#app").textContent = goHello();


// Load the webassembly and the executor
import '../../bin/wasm_exec'
import wasmBytes from '../../bin/spaghetti.wasm'

//Load the styling
import './spaghetti.css'

let _editor_instances = 0;

export class Editor {
    
    /** DOM Container of the canvas */
    container;
    canvas;
    #instance;

    constructor(options = {}) {
        this.container = options.container ?? document.body;
        this.#instance = _editor_instances;
    }

    /**
     * Starts the loader
     */
    async run() {
        //Create the canvas
        this.#createCanvas();

        //Load the go module
		const go = new Go();
        const result = await WebAssembly.instantiate(wasmBytes, go.importObject);

        //Run the module
        go.argv = [ `.spaghetti-instance-${this.#instance}` ];
        await go.run(result.instance);
    }

    #createCanvas() {
        if (this.canvas != null) return this.canvas;
        this.canvas = document.createElement('canvas');
        this.container.appendChild(this.canvas);
        this.canvas.classList.add('spaghetti-canvas');
        this.canvas.classList.add(`spaghetti-instance-${this.#instance}`);
        this.canvas.setAttribute('oncontextmenu', 'return false;');
        return this.canvas;
    }
}
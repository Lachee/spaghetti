
// Load the webassembly and the executor
import '../../resources/bin/wasm_exec'
import wasm from '../../resources/bin/spaghetti.wasm'

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

        this.#injectGlobals();
    }

    /**
     * Starts the loader
     */
    async run() {
        //Create the canvas
        this.#createCanvas();

        //Load the go module
		const go = new Go();
        let result = null;
        if (typeof(wasm) === 'string') {
            this.log('streaming assembly from ', wasm);
            result = await WebAssembly.instantiateStreaming(fetch(wasm), go.importObject);
        } else {
            this.log('loading raw assembly byes');
            result = await WebAssembly.instantiate(wasm, go.importObject);
        }

        //Run the module
        go.argv = [ `.spaghetti-instance-${this.#instance}` ];
        await go.run(result.instance);
    }

    /** injects our own runner to the import objects */
    #injectGlobals() {
        const decoder = new TextDecoder("utf-8");
        let outputBuffer = "";

        fs.writeSync = (fd, buf) => {
            outputBuffer += decoder.decode(buf);
            const nl = outputBuffer.lastIndexOf("\n");
            if (nl != -1) {
                this.log(outputBuffer.substr(0, nl));
                outputBuffer = outputBuffer.substr(nl + 1);
            }
            return buf.length;
        }
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

    log(message, ...params) {
        console.log('[spaghetti]', message, ...params);
    }
}

// Load the webassembly and the executor
import './spaghetti.css'
import '../../resources/bin/wasm_exec'
import wasm from '../../resources/bin/spaghetti.wasm'
import { fetchResource as fetchResource } from './resource';
import earcut from 'earcut';



let _editor_instances = 0;

export class Editor {
    
    /** DOM Container of the canvas */
    container;
    canvas;

    /** @type {Boolean} Indicating if the downloaded resources get stored in the resource collection */
    cacheResources = false;
    resources = { };

    #instance;

    #panic = false;

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

        // Hook into go so we can get the exit code
        const exit = go.exit;
        go.exit = (code) => {
            exit(code);
            go.exitCode = code;
        } 

        //Run the module
        go.argv = [ `.spaghetti-instance-${this.#instance}` ];
        await go.run(result.instance);
        
        //Show the panic
        if (this.#panic !== false)
            this.#displayPanic();
        
        //Return the exit code
        return go.exitCode;
    }

    /** injects our own runner to the import objects */
    #injectGlobals() {
        const decoder = new TextDecoder("utf-8");
        let outputBuffer = "";

        fs.writeSync = (fd, buf) => {
            outputBuffer += decoder.decode(buf);
            const nl = outputBuffer.lastIndexOf("\n");
            if (nl != -1) {

                const msg = outputBuffer.substr(0, nl);
                if (msg.startsWith("panic:")) this.#panic = [];
                
                if (this.#panic !== false) {
                    this.#panic.push(msg);
                } else {
                    this.log(msg);
                }

                // Start the new buffer
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
        this.canvas.editor = this;
        return this.canvas;
    }
    
    log(message, ...params) {
        console.log('[spaghetti]', message, ...params);
    }
    warn(message, ...params) {
        console.warn('[spaghetti]', message, ...params);
    }
    error(message, ...params) {
        console.error('[spaghetti]', message, ...params);
    }
    #displayPanic() {
        const panic = this.#panic.join('\n');
        this.error(panic);

        const panicBox = document.createElement('div');
        panicBox.classList.add("spaghetti-panic");
        panicBox.innerText = panic;

        //Append and hide the container
        this.container.appendChild(panicBox);
        this.canvas.style.display = 'none';
    }


    /** Fetches the given resource or url .
    */
    async fetchResource(resource) {
        console.groupCollapsed(`Loading resource "${resource}"`);
        try {
            if (this.cacheResources && this.resources[resource]) 
                return this.resources[resource];

            const result = await fetchResource(resource);
            if (this.cacheResources) this.resources[resource] = result;
            return result;
        } finally {

            console.groupEnd("Resource " + resource);
        }
    }

}

// Load the webassembly and the executor
import './spaghetti.css'
import '../../resources/bin/wasm_exec'
import wasm from '../../resources/bin/spaghetti.wasm'


export const ResourceProtocol = 'resource://';
export const Resources = {};
const context = require.context('../../resources/', true, /.*$/, 'sync');
context.keys()
        .forEach(path => {
            const name = path.substr(2);
            if (!name.startsWith('bin/')) {
                const m =  context(path);
                Resources[name] = m.default || m;
            }
        });


let _editor_instances = 0;

export class Editor {
    
    /** DOM Container of the canvas */
    container;
    canvas;

    /** @type {String} URL to the resource folder when it needs to look up. */
    resourceURL = '/resources';
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
     * @return {Promise<Uint8Array|Image>}
    */
    async fetchResource(resource) {        
        this.log('resource', 'fetch', resource);

        // prepare results
        let url = resource;
        let results = null;

        // Its a resource, so load from there
        if (resource.startsWith(ResourceProtocol)) {
            resource = resource.substr(ResourceProtocol.length);                
            const module = this.resources[resource] || Resources[resource];
            
            // Return directly
            if (module instanceof Uint8Array || module instanceof Image) {
                results = module;
            } else if (module instanceof ArrayBuffer) {
                results = new Uint8Array(module);
            } else if (module !== null) {
               // If its a data url, then download the image.
                // otherwise we need to update our resource URL to it.
                if (typeof(module) === 'string') {     
                    if (module.startsWith('data:image')) {
                        this.log('resource', 'decode image'); 
                        results = await this.#loadImage(module);
                    } else if (module.startsWith('data:;')) {
                        this.log('resource', 'decode data'); 
                        const enc = new TextEncoder(); 
                        results = enc.encode(atob(module.substr(13)));
                    } else {
                        url = module;
                    }
                }
            }
        }
        
        // We have no results yet, so lets just download the resource
        if (results == null) {
            results = await this.downloadResource(url);
        }

        return results;
    }

    /** Downloads the URL 
     * @return {Promise<Uint8Array|Image>} the downloaded data or image.
    */
    async downloadResource(url) {
        this.log('resource', 'download', url); 
        const response = await fetch(url);
        const contentType = response.headers.get('content-type');

        if (contentType.startsWith('image/')) {
        
            // Convert the image data
            // this.log('resource', 'image file', url);
            const blob = await response.blob();
            return await this.#loadImage(blob, contentType);
        
        } else {

            // Convert the binary data
            // this.log('resource', 'binary file', url);
            const buff = await response.arrayBuffer();
            return new Uint8Array(buff);
        }
    }

    /** Loads a new image from the url and waits for it to be done.
     * @return {Promise<Image>} the image promise
     */
    async #loadImage(data, mime = 'image/png') {
        return await new Promise((resolve, reject) => {
            const img = new Image();
            img.setAttribute('crossOrigin', 'anonymous');
            img.onload = function() {
                resolve(img);
            }
            img.onerror = function(message) {
                reject(message);
            }

            if (typeof(data) === 'string') {                // Data is a URI so we can just use it directly
                img.src = data;
            } else if (data instanceof Uint8Array) {        // Data is an array of bytes so it needs converting
                const encoded = btoa(String.fromCharCode.apply(null, ascii));
                img.src = `data:${mime};base64,${encoded}`;
            } else if (data instanceof Blob) {              // Data is a blob
                img.src = URL.createObjectURL(data);
            }
        });
    }    
}
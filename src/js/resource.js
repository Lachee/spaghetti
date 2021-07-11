import opentype from 'opentype.js';
import { Font } from './font';

// Load the resources
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

// List of loaders
const loaders = {
    'text/plain':   stringLoader,
    'image/png':    imageLoader,
    'image/jpeg':   imageLoader,
    'image/gif':    imageLoader,
    'image/svg':    imageLoader,
    'font/ttf':     fontLoader,
}

async function fontLoader(buff) {

    // It's already a buffer so lets just return the parsed font
    if (buff instanceof ArrayBuffer) {
        const font = opentype.parse(buff);
        return new Font(font);
    }

    // It's a URI
    const blob = new Blob([buff]);
    const uri = URL.createObjectURL(blob);
    return await new Promise((resolve, reject) => {
        opentype.load(uri, function(err, font) {
            if (err) return reject(err);         
            resolve(new Font(font));
        });
    });
}
async function stringLoader(buff) {
    if (typeof buff === 'string') return buff;
    if (buff instanceof ArrayBuffer) buff = new Uint8Array(buff);
    return String.fromCharCode.apply(null, buff);
}
async function imageLoader(buff, resource) {
    return await new Promise((resolve, reject) => {
        const img = new Image();
        img.setAttribute('crossOrigin', 'anonymous');
        img.onload = function() { resolve(img); }
        img.onerror = function(message) { console.error('failed to load image', message); reject(message); }

        const blob = new Blob([buff]);
        const uri = URL.createObjectURL(blob);
        img.src = uri;
    });
}


/** Fetches the given resource or url .
 * @return {Promise<Uint8Array|Image>}
*/
export async function fetchResource(resource) {
    console.log('[resource]', 'fetch', resource);

    // Its a resource, so load from there
    if (resource.startsWith(ResourceProtocol)) {
        const name = resource.substr(ResourceProtocol.length);
        return await loadResource(name);
    }
    
    // We have no results yet, so lets just download the resource
    return await downloadResource(resource);
}

/** Loads the resource at the given path */
async function loadResource(resoucePath) {
    console.log('[resource]', 'load', resoucePath);
    const content = Resources[resoucePath];
    if (content == null) return null;
    
    const [ header, dataTail ] = content.split(';', 2);
    const [ _, mimeType ] = header.split(':', 2);
    const [ format, data ] = dataTail.split(',', 2);

    // If its a URI then immediately download it
    if (format === 'uri')
        return await downloadResource(data);

    let buffer  = Uint8Array.from(atob(data), c => c.charCodeAt(0))
    let result  = buffer;
    if (loaders[mimeType]) {
        result = await loaders[mimeType](buffer); 
    } else {
        console.warn('resource', 'missing loader for content type', contentType);
    }
    return [mimeType, result];
}

/** Downloads a resource */
async function downloadResource(url) {
    console.log('[resource]', 'download', url);
    const response = await fetch(url);
    const contentType = response.headers.get('content-type').split(';', 2)[0].trim();
    
    let data = await response.arrayBuffer();
    if (loaders[contentType])
        data = await loaders[contentType](data);
    
    //console.log('[resource]', 'downloaded', contentType, data);
    return [contentType, data];
}
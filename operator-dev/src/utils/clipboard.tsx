import copy from 'copy-to-clipboard';

export const useClipboard = () => {
    const doCopy = (msg: string) => {
        const copyPromise = new Promise(function (resolve, reject) {
            try {
                copy(msg)
                resolve(msg);
            } catch (error) {
                reject(error);
            }
        })
        return copyPromise
    }
    return doCopy
}
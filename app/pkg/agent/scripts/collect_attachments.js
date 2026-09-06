(async () => {

    const selList = `div[data-automation-id="fileAttachment"] div[data-automation-id="promptOption"]`;
    const selUrl = `div[data-automation-id="copyUrl"]`;
    const selText = `div[data-automation-id="copyText"]`;
    const attribute = 'data-clipboard-text';
    const sleep = (delay) => scheduler.postTask(() => { }, { delay });

    const result = [];

    const container = [].slice.call(window.document.querySelectorAll(selList));
    for (let i = 0; i < container.length; i++) {
        const div = container[i];
        {
            const e = div.ownerDocument.createEvent('MouseEvents');
            e.initMouseEvent('contextmenu', true, true,
                div.ownerDocument.defaultView, 1, 0, 0, 0, 0, false,
                false, false, false, 2, null);
            div.dispatchEvent(e);
            await sleep(500);
        }

        const url = window.document.querySelector(selUrl).getAttribute(attribute);
        const text = window.document.querySelector(selText).getAttribute(attribute);

        result.push({ Title: text, Url: url });

        {
            window.dispatchEvent(
                new KeyboardEvent("keydown", {
                    altKey: false,
                    code: "Escape",
                    ctrlKey: false,
                    isComposing: false,
                    key: "Escape",
                    location: 0,
                    metaKey: false,
                    repeat: false,
                    shiftKey: false,
                    which: 27,
                    charCode: 0,
                    keyCode: 27,
                })
            );
            await sleep(500);
        }
    }

    return result;
})();
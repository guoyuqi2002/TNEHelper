(() => {
    const list = [].slice.call(window.document.querySelectorAll(`div[data-automation-id="multiViewListDetailItem"]`));
    return list.map((item, idx) => {
        item.id = `TNE_HELPER_TMP_ID_${idx}`;

        return item.querySelectorAll(`svg.wd-icon`).length !== 0;
    })
})();
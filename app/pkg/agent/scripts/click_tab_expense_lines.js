(() => {
    const div = [].slice.call(window.document.querySelectorAll(`div[data-automation-id="tabLabel"]`)).find(x => x.innerText === "Expense Lines");
    div.click();
})();
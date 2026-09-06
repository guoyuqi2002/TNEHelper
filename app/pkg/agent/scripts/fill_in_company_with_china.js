(() => {
    const labels = [].slice.call(window.document.getElementsByTagName("label"));
    const labelCompany = labels.find(l => l.innerText === "Company");
    const inputCompany = labelCompany.parentNode.parentNode.getElementsByTagName("input")[0];
    inputCompany.click();


    const labelDocumentNumber = labels.find(l => l.innerText === "Document Number");
    const inputDocumentNumber = labelDocumentNumber.parentNode.parentNode.getElementsByTagName("input")[0];
    inputDocumentNumber.click();
    inputDocumentNumber.value = "${DocumentNumber}";

    inputDocumentNumber.dispatchEvent(new Event('input', { bubbles: true }));
    inputDocumentNumber.dispatchEvent(new Event('change', { bubbles: true }));
})();
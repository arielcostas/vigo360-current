let workId = window.location.pathname.split("/")[3]
let form = document.getElementById("attachment-upload")
let imageStatus = document.getElementById("attachment-status")

/**
 * @typedef {{title: string, filename: string}} Attachment
 * @returns {Promise<Attachment[]>}
 */
async function listAttachments() {
    let resp = await fetch('/admin/async/attachments?trabajo=' + workId)
    let body = await resp.json()

    /**
     * @type {Attachment[]} attachments
     */
    let attachments = []
    if (resp.status === 200) {
        attachments = body
    }
    return attachments
}

async function deleteAttachment(n) {
    let resp = await fetch("/admin/async/attachments?id=" + n, {
        method: "DELETE"
    })
    let body = await resp.json()

    if (resp.status === 200) {
        loadAttachments()
        imageStatus.innerText = `Adjunto borrado`
        setTimeout(() => {
            imageStatus.innerText = ""
        }, 2000)
    } else {
        imageStatus.innerText = body.errorMsg
    }
}

function loadAttachments() {
    let attachmentList = document.getElementById("attachment-list")
    imageStatus.innerHTML = "<b>Cargando datos...</b>"
    let attachments = listAttachments();
    attachments.then((items) => {
        if (items.length < 1) {
            imageStatus.innerHTML = "<b>Ningún adjunto añadido</b>"
            return;
        }

        imageStatus.innerHTML = ""
        attachmentList.innerHTML = ""

        items.forEach(i => {
            let li = document.createElement("li")
            li.setAttribute("class", "attachment")
            li.innerText = `${i.title} (${i.filename})`
            attachmentList.appendChild(li)
        })
    })

}

loadAttachments()

form.addEventListener("submit", async (e) => {
    e.preventDefault()
    let titulo = form.getElementById("attachment-title").value
    let archivo = form.getElementById("attachment-file").files[0]

    fd = new FormData()
    fd.append("titulo", titulo)
    fd.append("trabajo", workId)
    fd.append("file", archivo)

    if (archivo.size > 157_286_400) {
        imageStatus.innerText = `El archivo seleccionado es demasiado pesado`
        return
    }

    resp = await fetch('/admin/async/attachments', {
        method: "POST",
        body: fd
    })
    if (resp.status !== 200) {
        body = await resp.json()
        imageStatus.innerText = body.errorMsg
    } else {
        imageStatus.innerText = `Archivo guardado exitosamente`
        setTimeout(() => {
            imageStatus.innerText = ""
        }, 2000)
    }

    loadAttachments()
})
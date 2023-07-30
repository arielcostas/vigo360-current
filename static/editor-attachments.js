let workId = window.location.pathname.split("/")[3]
let form = document.getElementById("attachment-upload")
let imageStatus = document.getElementById("attachment-status")

/**
 * @typedef {{id: int, title: string, filename: string}} Attachment
 * @returns {Promise<Attachment[]>}
 */
async function listAttachments() {
    let resp = await fetch('/admin/async/attachments?trabajo=' + workId)
    let body = await resp.json()

    /**
     * @type {Attachment[]} attachments
     */
    let attachments = []
    if (resp.status == 200) {
        attachments = body
    }
    return attachments
}

async function deleteAttachment(n) {
    if (!window.confirm(`¿Seguro que desea eliminarlo? Esta acción es irreversible`)) {
        return;
    }

    let resp = await fetch("/admin/async/attachments?id=" + n, {
        method: "DELETE"
    })

    if (resp.status == 204) {
        loadAttachments()
        imageStatus.innerText = `Adjunto borrado`
        setTimeout(() => {
            imageStatus.innerText = ""
        }, 2000)
    } else {
        let body = await resp.json()
        imageStatus.innerText = body["error"]
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
            li.innerHTML = `
                <a onclick="deleteAttachment(${i.id})"><b>Borrar</b></a> &mdash; ${i.title}
                (<a href="/static/papers/${i.filename}" target="_blank">${i.filename}</a>)`
            attachmentList.appendChild(li)
        })
    })

}

loadAttachments()

form.addEventListener("submit", async (e) => {
    e.preventDefault()
    let titulo = document.getElementById("attachment-title").value
    let archivo = document.getElementById("attachment-file").files[0]

    fd = new FormData()
    fd.append("titulo", titulo)
    fd.append("trabajo", workId)
    fd.append("file", archivo)

    if (archivo.size > 157_286_400) {
        imageStatus.innerText = `El archivo seleccionado es demasiado pesado`
        return
    }

    let resp = await fetch('/admin/async/attachments', {
        method: "POST",
        body: fd
    })
    if (resp.status != 201) {
        let body = await resp.json()
        imageStatus.innerText = body["error"]
    } else {
        imageStatus.innerText = `Archivo guardado exitosamente`
        setTimeout(() => {
            imageStatus.innerText = ""
        }, 2000)
        form.reset()
    }

    loadAttachments()
})
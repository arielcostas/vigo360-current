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

async function eliminarFotografia(n) {
    let resp = await fetch("/admin/async/a?foto=" + n, {
        method: "DELETE"
    })
    let body = await resp.json()

    if (resp.status === 200) {
        loadAttachments()
        imageStatus.innerText = `Imagen borrada`
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
    archivo = document.getElementById("archivo")
    fd = new FormData()
    fd.append("foto", archivo.files[0])
    fd.append("articulo", workId)

    if (archivo.files[0].size > 20971520) {
        imageStatus.innerText = `La imagen seleccionada es demasiado pesada`
        return
    }

    resp = await fetch('/admin/async/fotosExtra', {
        method: "POST",
        body: fd
    })
    if (resp.status !== 200) {
        body = await resp.json()
        imageStatus.innerText = body.errorMsg
    } else {
        imageStatus.innerText = `Imagen guardada exitosamente`
        setTimeout(() => {
            imageStatus.innerText = ""
        }, 2000)
    }

    loadAttachments()
})
document.getElementById("portada").addEventListener('change', e => {
	let selected = document.getElementById("portada").files[0];

	let reader = new FileReader();
	reader.addEventListener("load", () => {
		let imgTag = document.getElementById("portada_actual");
		imgTag.src = reader.result;
	});
	reader.readAsDataURL(selected);
});

function resize(){
	var login = document.getElementById("login");
	var reg = document.getElementById("register");
	login.style.marginTop = (window.innerHeight - login.scrollHeight) * 0.5 + "px";
	reg.style.opacity = 0;
	reg.style.display = "table";
	reg.style.marginTop = (window.innerHeight - reg.scrollHeight) * 0.5 + "px";
	reg.style.display = "none";
	reg.style.opacity = 100;
}
function hideLogin() {
	document.getElementById("login").style.display = "none"
	document.getElementById("register").style.display = "table"
	document.getElementById("regWrong").style.display = "block"
	document.getElementById("autoWrong").style.display = "none"
}
function showLogin() {
	document.getElementById("login").style.display = "table"
	document.getElementById("register").style.display = "none"
	document.getElementById("regWrong").style.display = "none"
	document.getElementById("autoWrong").style.display = "block"
}
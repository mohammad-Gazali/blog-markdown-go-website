const themeToggleDarkIcons = document.querySelectorAll(
	".theme-toggle-dark-icon"
);
const themeToggleLightIcons = document.querySelectorAll(
	".theme-toggle-light-icon"
);

// Change the icons inside the button based on previous settings
if (
	localStorage.getItem("color-theme") === "dark" ||
	(!("color-theme" in localStorage) &&
		window.matchMedia("(prefers-color-scheme: dark)").matches)
) {
	themeToggleLightIcons.forEach((t) => t.classList.remove("hidden"));
	document.documentElement.classList.add("dark");
} else {
	themeToggleDarkIcons.forEach((t) => t.classList.remove("hidden"));
	document.documentElement.classList.remove("dark");
}

const themeToggleBtns = document.querySelectorAll(".theme-toggle");

themeToggleBtns.forEach((btn) =>
	btn.addEventListener("click", function () {
		// toggle icons inside buttons
		themeToggleDarkIcons.forEach((t) => t.classList.toggle("hidden"));
		themeToggleLightIcons.forEach((t) => t.classList.toggle("hidden"));

		// if set via local storage previously
		if (localStorage.getItem("color-theme")) {
			if (localStorage.getItem("color-theme") === "light") {
				document.documentElement.classList.add("dark");
				localStorage.setItem("color-theme", "dark");
			} else {
				document.documentElement.classList.remove("dark");
				localStorage.setItem("color-theme", "light");
			}

			// if NOT set via local storage previously
		} else {
			if (document.documentElement.classList.contains("dark")) {
				document.documentElement.classList.remove("dark");
				localStorage.setItem("color-theme", "light");
			} else {
				document.documentElement.classList.add("dark");
				localStorage.setItem("color-theme", "dark");
			}
		}
	})
);

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.{templ,html}"],
    theme: {
        extend: {},
    },
    plugins: [require("@tailwindcss/forms")],
};

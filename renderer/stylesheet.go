package renderer

const stylesheet = `

body {
	font-family: "Verdana", Verdana, Arial, sans-serif;
	font-size: 12px;
}

a {
	text-decoration: none;
}

#toc {
	position: fixed;
	height: 98%;
	width: 400px;
	overflow-y: auto;
	top: 0;
	left: 0;
}

#toc ul {
	list-style: none;
	padding-left: 10px;
	padding-right: 10px;
}

li {
	min-height: 20px;
}

#main {
	margin-left: 370px;
  }

.entity {
	border: 1px solid black;
	padding: 10px;
	margin: 25px 0px 10px 30px;
}

.title {
	font-weight: bold;
	font-size: 110%;
}

.description {
	font-style: italic;
	padding: 10px 0;
}

.label {
	background: #e4b9c0;
	padding: 4px;
	border-radius: 6px;
	text-transform: uppercase;
	font-weight: normal;
	margin: 0 5px;
}

.fields li {
	padding-bottom: 10px;
}

`

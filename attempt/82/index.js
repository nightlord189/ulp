const express = require('express');
const path = require('path');

const port = 8080;


const app = express();

app.use(express.static(path.join(__dirname, 'public')))

app.get('/user/add', (req, res) => {
    res.send({ message: 'user successfully added!' });
});

app.get('/index', (req, res) => {
    res.sendFile(`${__dirname}/public/index.html`);
});

app.listen(port, () => {
    console.log(`Application listening on port ${port}!`);
});

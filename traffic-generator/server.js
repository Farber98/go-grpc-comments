
  
const express = require('express');
const router = express.Router();

const app = express();

// Settings 
app.set('port', process.env.PORT || 3000);

// Middlewares
app.use(express.json());

// Routes
app.use('/', router);

router.post('/', async (req, res) => {
    const { name, email, comment } = req.body;
    console.log({recibo: 'ACK', name:name, email:email, comment:comment});
    res.json({recibo: 'ACK', name:name, email:email, comment:comment});
});

// Starting the server
app.listen(app.get('port'), () => {
  console.log(`Server on port ${app.get('port')}`);
});
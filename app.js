const express = require("express");
const morgan = require("morgan");
const app = express();
const cors = require("cors");

// Middlewares
app.use(morgan("dev"));
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

//Routes
app.use("/users", require("./routes/users"));

// Start the Server
const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`Server Start, ${port}port`);
});

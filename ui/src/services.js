import axios from "axios";
const dispatchServer = (action) => {

    axios
    .post("/dispatch", action) // GET request to URL /hello
    .then(resp => console.log(resp.data)) // save response to state
    .catch(err => console.log(err)); // catch error
}


export default {
    dispatchServer
}
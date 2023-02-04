import axios from 'axios'
import React, {useState} from "react";
import { createBrowserHistory } from "history";

const Signup =()=>{
 const [name, setName] = useState("") 
 const [email, setEmail] = useState("")
 const [password, setPassword] = useState("")
 const [error, setError] = useState("")
 const [signedUp, setSignedUp] = useState(false)

 function inputName(e){
   e.preventDefault()
   setName(e.target.value)
 }

 function inputEmail(e){
   e.preventDefault()
   setEmail(e.target.value)
 }

 function inputPassword(e){
  e.preventDefault()
  setPassword(e.target.value)
}


function handleSubmit(e){
  var re = {
    'capital' : /[A-Z]/,
    'digit'   : /[0-9]/,
    'number'  : /\d/,
    'special' : /[ `!@#$%^&*()_+\-=\]{};':"\\|,.<>?~]/,
    'eMail'   : /(?!.*\.{2})^([a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+(\.[a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+)*|"((([\t]*\r\n)?[\t]+)?([\x01-\x08\x0b\x0c\x0e-\x1f\x7f\x21\x23-\x5b\x5d-\x7e\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|\\[\x01-\x09\x0b\x0c\x0d-\x7f\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))*(([\t]*\r\n)?[\t]+)?")@(([a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.)+([a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.?$/i,
};

const checkName = name
const checkEmail = email
const checkPassword = password
let passwordLength = password.length

  let axiosConfig = {
    headers: {
        'Content-Type': 'application/json;charset=UTF-8',
        "Access-Control-Allow-Origin": "*",
    }
  };

  if(passwordLength<8){
    setError("Password is too short")
   }
   else if(re.eMail.test(checkEmail)===false){
     setError("Email address is not valid")
   }
   else if(re.capital.test(checkPassword)===false){
     setError("Your password needs at least one capital letter")
   }
   else if(re.number.test(checkPassword)===false){
     setError("Your password needs at least one number")
   }
   else if(re.special.test(checkPassword)===false){
     setError("Your password needs at least one special character")
   }
else{
  const history = createBrowserHistory();
  history.push({
      pathname: `/account`
    });
 axios.post('http://localhost:8000/api/postform', {
  name,
  email,
  password,
})
.then(function (response) {
  console.log(response);
})
.catch(function (error) {
  console.log(error);
});
}
setSignedUp(true)
}
      return(
        <div>
         {console.log(name)}
         {console.log(email)}
         {console.log(password)}
         {console.log(error)}
          <form onSubmit={handleSubmit}>
          <h1 className="auth-heading">Sign up</h1>
          <input type="text" placeholder="FullName" name="query" required={true} onChange={inputName} />
          <input type="text" placeholder="Email" name="query" required={true} onChange={inputEmail} />
          <input type="text" placeholder="Password" name="query" required={true} onChange={inputPassword} />
          <input type="submit" value="Create Account" className="submit"/>
          {console.log(signedUp)}
          </form>
        </div>
      );
  }

export default Signup;
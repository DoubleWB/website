export default class SignatureService {
    static myInstance = null;
    //Create new singleton
    static getInstance() {
        if (SignatureService.myInstance == null) {
            SignatureService.myInstance =
                new SignatureService();
        }
        return this.myInstance;
    }

    //Fetch all signatures
    getSignatures = () =>
        fetch(`https://www.doublewb.xyz/api/signs`, {
            method: 'GET',
        })
            .then(response => response.json())

    //Send new signature to the server
    createSignature = signature =>
        fetch(`https://www.doublewb.xyz/api/sign`, {
            method: 'POST',
            body: JSON.stringify(signature),
            headers: {
                'content-type': 'application/json'
            }
        })
            .then(response => response.json())
}
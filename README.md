## authchain

**Authentic Chains**: Self-validating chains of entries with authenticity guarantees based on digital signatures.

Contraction of Authentic or Authored Chains.
This is an abstract distributed data structure to represent and store a sequence or sequences of digital messages such that:

A chain is permanently connected to a single public identifyer.
The public identifyers are connected to asymetric key pairs such that the private key is secured and available to the owner of the identity for signing (authenticating) new messages to be added to the chain, or decrypting private messages.
A creation message is authenticated to the identity by generating a signature of the creation message using the private key.

All messages are stored in the distributed database indexed by its signature.
The identity's public key can be used at any time to validate the authenticity of the message id (storage index) with respect to the message contents.

A new message can always be added to a given point in an existing chain.
Each new message is chained to the message before by its message id by including it in the new message and then generating the signature for the new message (using the identity's private key to validate the message is authentic).
In this way, any message id can be used to authenticate the entire message chain from the new message back to the creation message.

The history of every message is a unique chain back to the header, but it is also valid to add new messages again and again to the same point, even to the creation message, of a given chain.
Therefore, a nessage id does not uniquely specify the next message, or even whether any next message exists.
Applications can use other means to identify the first, or other designated singular next message added to a chain.

An object oriented design will have the following types of object:

* Digital identity: Active identities will have access to both the public and private keys. Ideally the private keys are held in a secure computing device that can sign and encrypt messages using the private key it stores, but never making any copies of that private key (Note 1 ref). This device would generate key pairs and export the public key for storage in the directory service where it will be associated with both a unique ID and other more human centered means of looking up a person's or trusted system's identity and public keys. 
* SecureDevice: Represents a service that will securely store the private key and use that key to encrypt messages or generate digital signatures for them.
* AuthChain: These are created with a creation message and an active identity. This will create a header record that is then signed by the active identity. This signature is the message ID used to index the message in a store. Given an existing AuthChain and an active identity, a new message can be added and signed generating a new message ID for the new state of the AuthChain.
* Message and Message ID: In principle AuthChains are agnostic about message formatting and encoding, however AuthChain applications need to agree. The reference implementation will probable use XML and/or other general data structures that can represent nested externalized data. Objects will need to be serialized into a cannonical format so that all instances of an application will serialize any object in the same way and generate/validate consistently. Stated in terms of invariants, if and only if two objects are the same (equal) they must serialize to the same binary message and 1) compute the same signature (message ID) with the private key or 2) validate that the message and message ID against the public identity (key) of the chain. 

## API Summary

### Identity

* Identity.new( secureDevice, message ) -> URL The message will need to identify this user publicly so it needs to include all attributes that will be used by applications to identify and look up public identities. This action will ask the secureDevice to generate a key pair as will as calling the directory service to store a new public identity. The URL will specify the new activeIdentity stored in the secureDevice and can be used with SecureDevice.new( URL )
* Identity.load( URL|Name|FileDescriptor ) -> activeIdentity: Opens an existing identity. Systems settings can select specific secureDevices to try, or it may be specified in a file by name, etc. Alternatively, the secureDevice is contained in the URL.

### SecureDevice

* SecureDevice.new( URL ) -> secDev A URL like: local:path, localdevice:path or FQDNofHost:path. Cases are 1) a local service, 2) A local device (maybe USB or Wifi/Bluetook key fob or 3) a service on the nework. This call opens a connection to the device and returns an object representing the open service or an error if it can't connect.
* secDev.generate( type ) -> localIdentityID  The key pair are stored in the device. The action creates a new localIdentityID in the secure device and it can be openned with new and a URL referencing the secure device and this ID.
* secDev.publicKey( [localIdentity] ) -> publicKey  If the argument is not supplied, the URL must have the ID in it (i.e. you openned an existing identity in the device).
* secDev.sign( message ) -> signature  Computes the signature of a message.
* secDev.encrypt( message ) -> messageEncrypted  Encrypt with private key
* secDev.decrypt( messageEncrypted, [public key] ) -> messageCleartext Decrypt with public key
* secDev.clone( publicKey, [localIdentityID] ) -> cloneMessageEncrypted: Create a clone message containing a private key and encrypt it with the public key of the target secureDevice. For loadClone to work this public key must corespond to the localIdentityID of the secDev you call loadClone on.
* secDev.loadClone( cloneMessageEncrypted ) -> localIdentityID: The secDev object must have the identity of a private key that will decrypt the message (the target) and this method will return a new local ID for the loaded clone.
* secDev.valid?( message, messageID, [ publicKey ] ) Validates the message and ID, returns true when valid.

### AuthChain

* AuthChain.new( activeIdentity, headerMessage ) -> chain: Creates a new AuthChain with no entries (just the header/creation message that will specify what this chain is about, etc. Returns a chain object with no entries.
* AuthChain.get( messageID ) -> chainEntry: Lookup and load a specific chain state. If the ID is a chainID, the chain is empty, it has no messages added to it.
* chain.getID -> messageID: Returns ID of the current (last) entry in the chain.
* chain.getChainID -> chainID: Returns ID of the zero entry (header/creation message) of the chain.
* chain.getMessage -> message: Returns the deserialized message as an object. It will be the header message when the chain is empty
* chain.addEntry( message ) -> newChain: Add message to the chain and returns the resulting chain.

Note 1: Even though these key devices don't normally make copies of a private key, there could be a special case where it is authorized to make a clone of an active identity in a second secure computing device.



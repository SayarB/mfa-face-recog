from flask import Flask, request
from deepface import DeepFace
app = Flask(__name__)


@app.route('/')
def hello_world():
    return 'Hello World!'

# post route that takes an image and returns a json response
@app.route('/api/v1/face-recognition', methods=['POST'])
def face_recognition():
    # get the image from the request body
    image = request.files['image']
    # save image as a jpg file
    name = request.form['name']
    image.save('./uploads/'+name+'.jpg')
    # perform face recognition on the image
    ans  = DeepFace.verify('./uploads/'+name+'.jpg', './images/'+name+'.jpg', model_name='Facenet', detector_backend='mtcnn', enforce_detection=False)
    return {'status':'success', 'verified':dict(ans)['verified'], 'distance':dict(ans)['distance'], 'threshold':dict(ans)['threshold']}

# post route to register a new user
@app.route('/api/v1/register', methods=['POST'])
def register():
    # get the image from the request body
    image = request.files['image']
    # save image as a jpg file
    name = request.form['name']
    image.save('./images/'+name+'.jpg')
    
    return {'status':'success'}
    
if __name__ == '__main__':
    app.run(debug=True)
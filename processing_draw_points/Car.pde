//One car
class Car {

  int carIndex;
  float w,h;   // width and height
  float angle; 
  color carColor;
  ArrayList pointNodes;
  
  // Car Constructor
  Car(int _carIndex, float _x, float _y) {
    pointNodes = new ArrayList();
    carIndex = _carIndex;
    addPoint(_x, _y);
    carColor = getColor();
  }
  
  //Add a new point for this car
  void addPoint(float _x, float _y) {
    pointNodes.add(new PointNode(_x, _y));
  }
  
  //Display the cars movement
  void display() {
    ellipseMode(CENTER);  // Set ellipseMode to CENTER
    PointNode lastPoint;
    fill(carColor);
    stroke(carColor);
    if (pointNodes.size() > 0) {
       lastPoint = (PointNode) pointNodes.get(0);
       ellipse(lastPoint.x, lastPoint.y, 30, 30); 
     
       if (pointNodes.size() > 1) {
         for (int i = 1; i < pointNodes.size(); i++) { 
           PointNode p = (PointNode) pointNodes.get(i);
           ellipse(p.x, p.y, 30, 30);
           line(lastPoint.x, lastPoint.y, p.x, p.y);
           lastPoint = p; 
         }
       }
  
    } 
  }
  
  //Get this cars color
  color getColor() {
    float redColor = carIndex*random(20,254) % 254; 
    float blueColor = carIndex*random(20,254) % 254; 
    float greenColor = carIndex*random(20,254) % 254; 
    return color(redColor, blueColor, greenColor);
  }
  
}

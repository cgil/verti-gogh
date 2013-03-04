ArrayList cars;
int totalCars = 2;

 // setup() runs once - Config
void setup() { 
  size(800, 600);
  frameRate(1);
  
  cars = new ArrayList();
  
  for (int i = 0; i < totalCars; i++) {
    Car car = new Car(i, random(0,width-1), random(0,height-1));
    cars.add(car);
  }
}

 // draw() loops forever, until stopped
void draw() { 
  background(255);
  for(int i = 0; i < cars.size(); i++) {
     Car c = (Car) cars.get(i);
     c.addPoint(random(0,width-1), random(0,height-1));
     c.display();
  } 
}

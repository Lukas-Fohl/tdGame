from enum import Enum
from typing import Any

class blockType_Class(Enum):
    normal = 1,
    path = 2,
    pathStart = 3,
    pathEnd = 4,
    blocked = 5

class DirectionType_Class(Enum):
    up = 0,
    left = 1,
    down = 2,
    right = 3,
    end = 4,
    start = 5

class path_Class:
    pathDirs: []
    startX: int
    startY: int
    def __init__(self, startXIn: int, startYIn: int) -> None:
        self.pathDirs = []
        self.startX = startXIn
        self.startY = startYIn
        pass

class etk_class:
     positionX: float
     positionY: float
     path: path_Class
     directionsIdx: int
     directionsProg: float
     speedFactor: float
     def __init__(self, pathIn: path_Class, speedFactor = 0.2) -> None:
        self.positionX = pathIn.startX
        self.positionY = pathIn.startY
        self.path = pathIn
        self.directionsIdx = 0
        self.directionsProg = 0
        self.speedFactor = speedFactor
        pass
     def tick(self):
        if(round(self.directionsProg,2) >= 0.999):
            self.directionsProg = 0
            if(self.directionsIdx + 1 < len(self.path.pathDirs)):
                self.directionsIdx += 1
        match(self.path.pathDirs[self.directionsIdx]):
            case DirectionType_Class.start:
                if(self.directionsIdx + 1 < len(self.path.pathDirs)):
                    self.directionsIdx += 1
                self.directionsProg = 0
            case DirectionType_Class.end:
                self.directionsProg = 0
                return
            case DirectionType_Class.up:
                if(self.directionsProg + self.speedFactor > 1):
                    self.positionY += abs(1 - self.directionsProg)
                    self.directionsProg = 1
                else:
                    self.positionY += self.speedFactor
                    self.directionsProg += self.speedFactor
            case DirectionType_Class.down:
                if(self.directionsProg + self.speedFactor > 1):
                    self.positionY -= abs(1 - self.directionsProg)
                    self.directionsProg = 1
                else:
                    self.positionY -= self.speedFactor
                    self.directionsProg += self.speedFactor
            case DirectionType_Class.left:
                if(self.directionsProg + self.speedFactor > 1):
                    self.positionX -= abs(1 - self.directionsProg)
                    self.directionsProg = 1
                else:
                    self.positionX -= self.speedFactor
                    self.directionsProg += self.speedFactor
            case DirectionType_Class.right:
                if(self.directionsProg + self.speedFactor > 1):
                    self.positionX += abs(1 - self.directionsProg)
                    self.directionsProg = 1
                else:
                    self.positionX += self.speedFactor
                    self.directionsProg += self.speedFactor
        return

class block_Class:
    blockType: blockType_Class
    used = False
    placeable = True
    def __init__(self, blockTypeIn = blockType_Class.normal) -> None:
            self.blockType = blockTypeIn
            self.used = False
            match(self.blockType):
                case blockType_Class.normal:
                      self.placeable = True
                case blockType_Class.path:
                      self.placeable = False
                case blockType_Class.blocked:
                      self.placeable = False
            pass


class map:
    content: [block_Class]
    width = 0
    height = 0
    def __init__(self, width, height) -> None:
        self.content = []
        self.width = width
        self.height = height
        for x in range(width):
            tempLine = []
            for y in range(height):
                tempLine.append(block_Class())
            self.content.append(tempLine)

    def printMap(self):
        for y in range(self.height):
            tempLine = ""
            for x in range(self.width):
                match(self.content[x][y].blockType):
                    case blockType_Class.normal:
                          tempLine += "~~"
                    case blockType_Class.path:
                          tempLine += ".."
                    case blockType_Class.pathStart:
                          tempLine += "S."
                    case blockType_Class.pathEnd:
                          tempLine += "E."
                    case blockType_Class.blocked:
                          tempLine += "XX"
                
            print(tempLine)
    def applayPath(self, pathIn: path_Class):
        tempX = pathIn.startX
        tempY = pathIn.startY
        self.content[pathIn.startX][pathIn.startY].blockType = blockType_Class.pathStart
        for i in pathIn.pathDirs:
            match(i):
                case DirectionType_Class.up:
                    tempY += 1
                case DirectionType_Class.down:
                    tempY -= 1
                case DirectionType_Class.right:
                    tempX += 1
                case DirectionType_Class.left:
                    tempX -= 1
                case DirectionType_Class.end:
                    self.content[tempX][tempY].blockType = blockType_Class.pathEnd
                    break
                case _:
                    break
            self.content[tempX][tempY].blockType = blockType_Class.path
        return

class tower_Class:
    dmg: int
    offMs: int
    positionX: int
    positionY: int
    def __init__(self, positionXIn: int, positionYIn: int, dmgIn = 1, offMsIn = 300) -> None:
        self.dmg = dmgIn
        self.offMs = offMsIn
        self.positionX = positionXIn
        self.positionY = positionYIn
        pass

def main():
    myMap = map(20,20)
    myPath = path_Class(0,0)

    myPath.pathDirs.append(DirectionType_Class.start)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.up)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.right)
    myPath.pathDirs.append(DirectionType_Class.end)

    myMap.applayPath(myPath)

    myETK = etk_class(myPath, speedFactor=0.2)

    for i in range(150):
        if(myETK.path.pathDirs[myETK.directionsIdx] != DirectionType_Class.end):
            myETK.tick()
            #print(str(round(myETK.positionX,2)) + " ; " + str(round(myETK.positionY,2)))
            myMap.content[int(round(myETK.positionX,2))][int(round(myETK.positionY,2))].blockType = blockType_Class.path
        else:
            myMap.content[int(round(myETK.positionX,2))][int(round(myETK.positionY,2))].blockType = blockType_Class.pathEnd
        continue
    myMap.content[myPath.startY][myPath.startX].blockType = blockType_Class.pathStart

    myMap.printMap()
    return


if __name__ == "__main__":
     main()


"""
PATH:
    start square
    store array of directions --> apply to map as used blocks
        0
    1       3
        2
"""
meow_protocol        = Proto("Meow", "Meow Protocol")

version              = ProtoField.uint8("meow.version", "Version", base.DEC, nil, 0xF0)
cuteness             = ProtoField.uint8("meow.cuteness", "Cuteness", base.DEC, nil, 0xF)
messageType          = ProtoField.uint8("meow.messagetype", "Message type", base.DEC)
animalType           = ProtoField.uint8("meow.animaltype", "Animal type", base.DEC)
breed                = ProtoField.uint16("meow.breed", "Breed", base.DEC)
nameLength           = ProtoField.uint64("meow.namelen", "Name length", base.DEC)
name                 = ProtoField.string("meow.name", "Name", base.UNICODE, "The name of the entity meowing or performing the action")

meow_protocol.fields = { version, cuteness, messageType, animalType, breed, nameLength, name } 
local values         = {
    message = {
        [1] = { "Meow" },
        [2] = { "Purr" },
        [3] = { "Scratch" },
        [4] = { "Bite" },
        [5] = { "Paw at you" },
        [6] = { "Growl" },
        [7] = { "Hiss" }
    },
    animal = {
        [1] = { "Cat" },
        [2] = { "Human, male" },
        [3] = { "Human, female" },
        [4] = { "Human, unspecified" }
    },
    breed = {
        -- Cats
        [1] = {
            [0] = { "Unspecified" },
            [1] = { "Calico" },
            [2] = { "White" },
            [3] = { "Siamese" },
        },
    },
    cuteness = {
        [1] = "Sender is not cute.",
        [2] = "Sender is cute",
        [3] = "Sender is very cute",
        [4] = "Sender is extremely cute",
        [5] = "Sender does not know if they are cute",
        [6] = "Sender does not want to disclose if they are cute or not",
        [7] = "Sender is ugly",
        [8] = "Sender is extremely ugly",
    }
}

local function get_message_type(type)
    print("d",type)
    return values.message[type][1] or "Unknown"
end

local function get_animal_type(animal)
    print("d",animal)
    local animal_row = values.animal[animal]
    if animal_row then return animal_row[1] or "Unknown" else return "Unknown" end
end

local function get_breed(animal, breed)
    local breeds = values.breed[animal] or {}
    if breeds == {} then return "Unknown" end
    if breeds[breed] then return breeds[breed][1] or "Unknown" else return "Unknown" end
end

local function get_cuteness(cuteness)
    return values.cuteness[cuteness] or "Unknown"
end

function meow_protocol.dissector(buffer, pinfo, tree)
    length = buffer:len()
    if length == 0 then return end

    pinfo.cols.protocol = meow_protocol.name

    local subtree = tree:add(meow_protocol, buffer(), "Meow Protocol")

    local byte0 = buffer(0, 1)
    local messageTypeBuf = buffer(1, 1)
    local animalTypeBuf = buffer(2, 1)
    local breedBuf = buffer(3, 2)
    
    local nameLenBuf = buffer(5, 1)
    local nameLen = nameLenBuf:uint()
    local nameBuf = buffer(6, nameLen)

    print(byte0, messageTypeBuf, animalTypeBuf, breedBuf)

    local cutenessVal = bit.band(byte0:uint(), 0x0F)

    subtree:add(version, byte0)
    subtree:add(cuteness, byte0):append_text(" (" .. get_cuteness(cutenessVal) .. ") ")
    subtree:add(messageType, messageTypeBuf):append_text(" (" .. get_message_type(messageTypeBuf:uint()) .. ")")
    subtree:add(animalType, animalTypeBuf):append_text(" (" .. get_animal_type(animalTypeBuf:uint()) .. ")")
    subtree:add(breed, breedBuf):append_text(" (" .. get_breed(animalTypeBuf:uint(), breedBuf:uint()) .. ")")
    subtree:add(nameLength, nameLenBuf)
    subtree:add(name, nameBuf)
end

local udp_port = DissectorTable.get("udp.port")
udp_port:add(32390, meow_protocol)

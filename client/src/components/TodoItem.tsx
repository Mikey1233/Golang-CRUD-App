import {
  Badge,
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  Spinner,
  Text,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";
import { LuFileEdit } from "react-icons/lu";
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
} from "@chakra-ui/react";
import React from "react";

async function mutateData(
  mutateType: string,
  method: string,
  todo: Todo,
  completed?: boolean
) {
  if (mutateType === "updateTodo" && todo.completed)
    return alert("todo is already completed");
  try {
    const res = await fetch(BASE_URL + `/todos/${todo._id}`, {
      method: method,
      body: JSON.stringify({ body: todo.body, completed: completed }),

      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!res.ok) {
      throw new Error( "something went wrong");
    }
    return "";
  } catch (err) {
    console.log(err);
  }
}

const TodoItem = ({ todo }: { todo: Todo }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const initialRef = React.useRef(null);
  const finalRef = React.useRef(null);
  const EditedList: Todo = {
    _id: todo._id,
    body: "",
    completed: false,
  };
  const color2 = useColorModeValue("yellow.500", "yellow.100");
  const color = useColorModeValue("green.500", "green.300");
  const queryClient = useQueryClient();

  /////update todo status
  const { mutate: updateTodo, isPending: isUpdating } = useMutation({
    mutationKey: ["updateTodo"],
    mutationFn: () => mutateData("updateTodo", "PATCH", todo, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });

  ////////delete todo
  const { mutate: deleteTodo } = useMutation({
    mutationKey: ["deleteTodo"],
    mutationFn: () => mutateData("deleteTodo","DELETE",todo),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });
  ///////////// edit todo list
  const { mutate: editTodo } = useMutation({
    mutationKey: ["editTodo"],
    mutationFn: () => mutateData("editTodo", "PATCH", EditedList),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });
  return (
    <>
      <Flex gap={2} alignItems={"center"}>
        <Flex
          flex={1}
          alignItems={"center"}
          border={"1px"}
          borderColor={"gray.600"}
          p={2}
          borderRadius={"lg"}
          justifyContent={"space-between"}
        >
          <Text
            color={todo.completed ? color : color2}
            textDecoration={todo.completed ? "line-through" : "none"}
          >
            {todo.body}
          </Text>
          {todo.completed && (
            <Badge ml="1" colorScheme="green">
              Done
            </Badge>
          )}
          {!todo.completed && (
            <Badge ml="1" colorScheme="yellow">
              In Progress
            </Badge>
          )}
        </Flex>
        <Flex gap={2} alignItems={"center"}>
          <Box
            color={"green.500"}
            cursor={"pointer"}
            onClick={() => updateTodo()}
          >
            {!isUpdating && <FaCheckCircle size={20} />}
            {isUpdating && <Spinner size={"sm"} />}
          </Box>
          <Box cursor={"pointer"} onClick={onOpen}>
            <LuFileEdit />
          </Box>
          <Box
            color={"red.500"}
            cursor={"pointer"}
            onClick={() => deleteTodo()}
          >
            <MdDelete size={25} />
          </Box>
        </Flex>
      </Flex>
      <Modal
        initialFocusRef={initialRef}
        finalFocusRef={finalRef}
        isOpen={isOpen}
        onClose={onClose}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Edit list item</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={6}>
            <FormControl>
              <FormLabel>🥳🥳</FormLabel>
              <Input
                ref={initialRef}
                placeholder={todo.body || "enter new list"}
              />
            </FormControl>
          </ModalBody>

          <ModalFooter>
            <Button
              colorScheme="blue"
              mr={3}
              onClick={() => {
                // console.log(initialRef?.current?.value)
                let inputElement =
                  initialRef.current as unknown as HTMLInputElement;
                if (inputElement && inputElement.value != "") {
                  EditedList.body = inputElement.value;
                }
                return editTodo();
              }}
            >
              Edit
            </Button>
            <Button onClick={onClose}>Cancel</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};

export default TodoItem;

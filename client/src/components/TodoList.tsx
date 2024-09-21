import { Flex, Spinner, Stack, Text} from "@chakra-ui/react";

import TodoItem from "./TodoItem";
// import { data } from "framer-motion/client";
import { useQuery } from "@tanstack/react-query";

export type Todo = {
    _id: number ;
    body : string;
    completed : boolean
}

const TodoList = () => {
//   const [isLoading, setIsLoading] = useState(false);
const {data:todos,isLoading} = useQuery<Todo[]>({
    queryKey: ['todos'],
    queryFn:  async()=>{
        try {
    const res = await fetch("http://localhost:5000/api/todos")
    const data = await res.json()
    if (!res.ok){
        throw new Error(data.error || "something went wrong")
    }
    return data || []
        }catch(err){
            console.log(err) 
        }
    }
 })
  return (
    <>
      <Text
        bgGradient="linear(to-l, #0b85f8, #00ffff)"
        bgClip="text"
        textAlign={"center"}
        fontSize="4xl"
        my={2}
        textTransform={"uppercase"}
        fontWeight="abold"
      >
        Today's task
      </Text>
      {isLoading && (
        <Flex justifyContent={"center"} my={4}>
          <Spinner size={"xl"} />
        </Flex>
      )}
      {!isLoading && todos?.length === 0 && (
        <Stack alignItems={"center"} gap="3">
          <Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
            All tasks completed! 🤞
          </Text>
          <img src="/go.png" alt="Go logo" width={70} height={70} />
        </Stack>
      )}
      <Stack gap={3}>
        {todos?.map((todo) => (
          <TodoItem key={todo._id} todo={todo} />
        ))}
      </Stack>
    </>
  );
};
export default TodoList;

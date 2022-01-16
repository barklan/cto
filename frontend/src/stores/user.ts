import { acceptHMRUpdate, defineStore } from 'pinia'
import axios from "axios"

export const useUserStore = defineStore('user', () => {
  /**
   * Current named of the user.
   */
  const name = ref('')
  const id = ref('')
  const project = ref('')
  const token = ref('')

  const projectColor = ref('')
  const projectPrettyTitle = ref('')

  /**
   *
   * @param name - new name to set
   */
  function setName(newName: string) {
    name.value = newName
  }

  function pickColor() {
    var baseColors = ['#7c2d12', '#831843', '#14532d', '#164e63', '#312e81', '#4c1d95', '#701a75']
    var randomColor = baseColors[Math.floor(
      Math.random() * baseColors.length)];

    if (randomColor == projectColor.value) {
      randomColor = pickColor()
    }
    return randomColor
  }
  function setProjectPrettyTitle(projectID: string) {
    let url = import.meta.env.VITE_PROTOCOL + "://" + import.meta.env.VITE_HOSTNAME +
      "/api/porter/project/" + projectID
    axios.get(url)
      .then((response) => {
        if (response.status == 200 && response.data.title) {
          projectPrettyTitle.value = response.data.title
        } else {
          projectPrettyTitle.value = "Error fetching project title."
        }
      })
      .catch((error) => console.log(error));
  }

  function setID(newID: string) {
    id.value = newID
  }
  function setProject(newProject: string) {
    project.value = newProject
    projectColor.value = pickColor()
    setProjectPrettyTitle(project.value)
  }
  function setToken(newToken: string) {
    token.value = newToken
  }

  return {
    setName,
    setID,
    setProject,
    setToken,
    name,
    id,
    project,
    projectColor,
    projectPrettyTitle,
    token,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))

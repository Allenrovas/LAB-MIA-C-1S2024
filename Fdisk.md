void comandoFdisk(char Path[512], char Name[16], int Size, char Unit[25], char Type[25], char Fit[25], char Delete[25], int Add){
    FILE *archivo = fopen(Path, "rb+");//abrimos el archivo
    if (archivo== NULL){
        cout << "Error: No se pudo abrir el archivo" << endl;
        fclose(archivo);
        return;
    } else {
        MBR disk;
        fseek(archivo, 0, SEEK_SET);//nos posicionamos al inicio del archivo
        fread(&disk, sizeof(MBR), 1, archivo);//leemos el archivo
        string eliminacion = Delete;
        if(Add == 0 && eliminacion == "0"){
            int TemporalDesplz = 1+ sizeof(MBR);
            Partition particionExtendida;
            int indiceParticion = 0;
            bool nombreRepetido = false;
            bool VerificarEspacio = false;
            string Nombre = Name;
            string Tipo = Type;
            if(disk.mbr_partition_1.part_size != 0){
                //if(disk.mbr_partition_1.part_type == 'e')particionExtendida = disk.mbr_partition_1;
                if(disk.mbr_partition_1.part_type=='e')particionExtendida = disk.mbr_partition_1;
                //if(disk.mbr_partition_1.part_name,Name.c_str())nombreRepetido = true;
                if(strcmp(disk.mbr_partition_1.part_name,Nombre.c_str())==0)nombreRepetido = true;
                TemporalDesplz += VerificarEspacio? 0:disk.mbr_partition_1.part_size+1;
            } else{
                indiceParticion = 1;
                VerificarEspacio = true;
            }
            if(disk.mbr_partition_2.part_size != 0){
                if(disk.mbr_partition_2.part_type=='e')particionExtendida = disk.mbr_partition_2;
                if(strcmp(disk.mbr_partition_2.part_name,Nombre.c_str())==0)nombreRepetido = true;
                TemporalDesplz += VerificarEspacio? 0:disk.mbr_partition_2.part_size+1;
            } else if(!VerificarEspacio){
                indiceParticion = 2;
                VerificarEspacio = true;

            }
            if(disk.mbr_partition_3.part_size != 0){
                if(disk.mbr_partition_3.part_type=='e')particionExtendida = disk.mbr_partition_3;
                if(strcmp(disk.mbr_partition_3.part_name,Nombre.c_str())==0)nombreRepetido = true;
                TemporalDesplz += VerificarEspacio? 0:disk.mbr_partition_3.part_size+1;
            } else if(!VerificarEspacio){
                indiceParticion = 3;
                VerificarEspacio = true;
            }
            if(disk.mbr_partition_4.part_size != 0){
                if(disk.mbr_partition_4.part_type=='e')particionExtendida = disk.mbr_partition_4;
                if(strcmp(disk.mbr_partition_4.part_name,Nombre.c_str())==0)nombreRepetido = true;
                TemporalDesplz += VerificarEspacio? 0:disk.mbr_partition_4.part_size+1;
            } else if(!VerificarEspacio){
                indiceParticion = 4;
                VerificarEspacio = true;
            }
            string tipo = Type;
            if (indiceParticion == 0 and tipo != "l"){
                cout << "Error: Ya existen cuatro particiones" << endl;
                fclose(archivo);
                return;
            }
            if(nombreRepetido){
                cout << "Error: Ya existe una particion con ese nombre" << endl;
                fclose(archivo);
                return;
            }
            if(particionExtendida.part_type=='e' && tipo == "e"){
                cout << "Error: Ya existe una particion extendida" << endl;
                fclose(archivo);
                return;
            }
            if(tipo != "l"){
                Partition particion;
                particion.part_type = Type[0];
                particion.part_fit = Fit[0];
                particion.part_start = TemporalDesplz;
                particion.part_size = Size * ( Unit[0] == 'k'? 1024: Unit[0] == 'm' ? 1024 * 1024 : 1);
                strcpy(particion.part_name, Name);
                if(TemporalDesplz + particion.part_size > disk.mbr_tamano){
                    cout << "Error: No hay espacio suficiente para crear la particion" << endl;
                    fclose(archivo);
                    return;
                }
                if(indiceParticion == 1){
                    disk.mbr_partition_1 = particion;
                } else if(indiceParticion == 2){
                    disk.mbr_partition_2 = particion;
                } else if(indiceParticion == 3){
                    disk.mbr_partition_3 = particion;
                } else if(indiceParticion == 4){
                    disk.mbr_partition_4 = particion;
                }
                fseek(archivo, 0, SEEK_SET);
                fwrite(&disk, sizeof(MBR), 1, archivo);
                fclose(archivo);
                if (Type[0] == 'e'){
                    cout << "Se creo la particion extendida" << endl;
                } else{
                    cout << "Se creo la particion primaria" << endl;
                }
            } else{
                EBR ebr;
                if (particionExtendida.part_type != 'e'){
                    cout << "Error: No existe una particion extendida" << endl;
                    fclose(archivo);
                    return;
                }
                TemporalDesplz = particionExtendida.part_start;
                do {
                    fseek(archivo, TemporalDesplz, SEEK_SET);
                    fread(&ebr, sizeof(EBR), 1, archivo);
                    if(ebr.part_size != 0){
                        if(ebr.part_name == Name){
                            nombreRepetido = true;
                            break;
                        }
                        TemporalDesplz += ebr.part_size + sizeof(EBR)+1;
                    }
                } while (ebr.part_size != 0);
                if(nombreRepetido){
                    cout << "Error: Ya existe una particion con ese nombre" << endl;
                    fclose(archivo);
                    return;
                }
                if(TemporalDesplz+ sizeof(ebr) + Size * ( Unit[0] == 'k'? 1024: Unit[0] == 'm' ? 1024 * 1024 : 1) > particionExtendida.part_start + particionExtendida.part_size){
                    cout << "Error: No hay espacio suficiente para crear la particion" << endl;
                    fclose(archivo);
                    return;
                }
                EBR ebrNuevo;
                cout << "TemporalDesplz: " << TemporalDesplz << endl;
                cout << "Size of EBR: " << sizeof(EBR) << endl;
                ebrNuevo.part_fit = Fit[0];
                ebrNuevo.part_start = TemporalDesplz + sizeof(EBR)+1;
                ebrNuevo.part_size = Size * ( Unit[0] == 'k'? 1024: Unit[0] == 'm' ? 1024 * 1024 : 1);
                ebrNuevo.part_next = TemporalDesplz + sizeof(EBR)+1 + ebrNuevo.part_size;
                strcpy(ebrNuevo.part_name, Name);
                fseek(archivo, TemporalDesplz, SEEK_SET);
                fwrite(&ebrNuevo, sizeof(EBR), 1, archivo);
                fclose(archivo);
                cout << "Se creo la particion logica con exito" << endl;
            }
        }else if (Add != 0){
            string Nombre = Name;
            int indiceParticion = 0;
            if(strcmp(disk.mbr_partition_1.part_name,Nombre.c_str())==0){
                indiceParticion = 1;
            }
            if(strcmp(disk.mbr_partition_2.part_name,Nombre.c_str())==0){
                indiceParticion = 2;
            }
            if(strcmp(disk.mbr_partition_3.part_name,Nombre.c_str())==0){
                indiceParticion = 3;
            }
            if(strcmp(disk.mbr_partition_4.part_name,Nombre.c_str())==0){
                indiceParticion = 4;
            }

            if(indiceParticion!=0){
                int size = 0;
                if(disk.mbr_partition_1.part_name!="~"){
                    size += disk.mbr_partition_1.part_size;
                }
                if(disk.mbr_partition_2.part_name!="~"){
                    size += disk.mbr_partition_2.part_size;
                }
                if(disk.mbr_partition_3.part_name!="~"){
                    size += disk.mbr_partition_3.part_size;
                }
                if(disk.mbr_partition_4.part_name!="~"){
                    size += disk.mbr_partition_4.part_size;
                }
                if(indiceParticion==1){
                    Partition particion = disk.mbr_partition_1;
                    particion = disk.mbr_partition_1;
                    size -= particion.part_size;
                    particion.part_size += Add * ( Unit[0] == 'k'? 1024: Unit[0] == 'm' ? 1024 * 1024 : 1);
                    if(particion.part_size <= 0){
                        cout << "Error: No se puede reducir la particion a un tamaño menor o igual a 0" << endl;
                        fclose(archivo);
                        return;
                    }
                    if(size + particion.part_size <= disk.mbr_tamano){
                        int desplazamiento = particion.part_start + particion.part_size;
                        //actualizar particiones
                        disk.mbr_partition_1 = particion;
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);

                        fseek(archivo, desplazamiento, SEEK_SET);
                        fwrite(&particion, sizeof(Partition), 1, archivo);
                        cout << "Se actualizo la particion con exito" << endl;
                        fclose(archivo);
                        return;
                    } else{
                        cout << "Error: No hay espacio suficiente para aumentar la particion" << endl;
                        fclose(archivo);
                        return;
                    }
                } else if (indiceParticion==2){
                    Partition particion = disk.mbr_partition_2;
                    size -= particion.part_size;
                    particion.part_size += Add * ( Unit[0] == 'k'? 1024: Unit[0] == 'm' ? 1024 * 1024 : 1);
                    if(particion.part_size <= 0){
                        cout << "Error: No se puede reducir la particion a un tamaño menor o igual a 0" << endl;
                        fclose(archivo);
                        return;
                    }
                    if(size + particion.part_size <= disk.mbr_tamano){
                        int desplazamiento = particion.part_start + particion.part_size;
                        //actualizar particiones
                        disk.mbr_partition_2 = particion;
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);

                        fseek(archivo, desplazamiento, SEEK_SET);
                        fwrite(&particion, sizeof(Partition), 1, archivo);
                        cout << "Se actualizo la particion con exito" << endl;
                        fclose(archivo);
                        return;
                    } else{
                        cout << "Error: No hay espacio suficiente para aumentar la particion" << endl;
                        fclose(archivo);
                        return;
                    }
                } else if (indiceParticion==3) {
                    Partition particion = disk.mbr_partition_3;
                    size -= particion.part_size;
                    particion.part_size += Add * (Unit[0] == 'k' ? 1024 : Unit[0] == 'm' ? 1024 * 1024 : 1);
                    if (particion.part_size <= 0) {
                        cout << "Error: No se puede reducir la particion a un tamaño menor o igual a 0" << endl;
                        fclose(archivo);
                        return;
                    }
                    if (size + particion.part_size <= disk.mbr_tamano) {
                        int desplazamiento = particion.part_start + particion.part_size;
                        //actualizar particiones
                        disk.mbr_partition_3 = particion;
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);

                        fseek(archivo, desplazamiento, SEEK_SET);
                        fwrite(&particion, sizeof(Partition), 1, archivo);
                        cout << "Se actualizo la particion con exito" << endl;
                        fclose(archivo);
                        return;
                    } else {
                        cout << "Error: No hay espacio suficiente para aumentar la particion" << endl;
                        fclose(archivo);
                        return;
                    }

                } else if (indiceParticion==4) {
                    Partition particion = disk.mbr_partition_4;
                    size -= particion.part_size;
                    particion.part_size += Add * (Unit[0] == 'k' ? 1024 : Unit[0] == 'm' ? 1024 * 1024 : 1);
                    if (particion.part_size <= 0) {
                        cout << "Error: No se puede reducir la particion a un tamaño menor o igual a 0" << endl;
                        fclose(archivo);
                        return;
                    }
                    if (size + particion.part_size <= disk.mbr_tamano) {
                        int desplazamiento = particion.part_start + particion.part_size;
                        //actualizar particiones
                        disk.mbr_partition_4 = particion;
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);

                        fseek(archivo, desplazamiento, SEEK_SET);
                        fwrite(&particion, sizeof(Partition), 1, archivo);
                        cout << "Se actualizo la particion con exito" << endl;
                        fclose(archivo);
                        return;
                    } else {
                        cout << "Error: No hay espacio suficiente para aumentar la particion" << endl;
                        fclose(archivo);
                        return;
                    }
                }
            } else{
                cout << "Error: No se encontro la particion, se buscara en las extendidas" << endl;
                //buscar particion extendida
                Partition particionExtendida;
                int indiceParticion = 0;

                if(disk.mbr_partition_1.part_type == 'e'){
                    particionExtendida = disk.mbr_partition_1;
                    indiceParticion = 1;
                } else if(disk.mbr_partition_2.part_type == 'e'){
                    particionExtendida = disk.mbr_partition_2;
                    indiceParticion = 2;
                } else if(disk.mbr_partition_3.part_type == 'e'){
                    particionExtendida = disk.mbr_partition_3;
                    indiceParticion = 3;
                } else if(disk.mbr_partition_4.part_type == 'e'){
                    particionExtendida = disk.mbr_partition_4;
                    indiceParticion = 4;
                }
                if(indiceParticion == 0){
                    cout << "Error: No se encontro una particion extendida" << endl;
                    fclose(archivo);
                    return;
                }
                //buscar particion logica
                EBR ebr;
                int TemporalDesplz = particionExtendida.part_start;
                int size = 0;
                do {
                    fseek(archivo, TemporalDesplz, SEEK_SET);
                    fread(&ebr, sizeof(EBR), 1, archivo);
                    if(ebr.part_size != 0){
                        size += ebr.part_size;
                        if(strcmp(ebr.part_name,Nombre.c_str())==0){
                            cout << "Se encontro la particion logica" << endl;
                            size -= ebr.part_size;
                            ebr.part_size += Add * (Unit[0] == 'k' ? 1024 : Unit[0] == 'm' ? 1024 * 1024 : 1);
                            if (ebr.part_size <= 0) {
                                cout << "Error: No se puede reducir la particion a un tamaño menor o igual a 0" << endl;
                                fclose(archivo);
                                return;
                            }
                            if (size + ebr.part_size <= disk.mbr_tamano) {
                                int desplazamiento = ebr.part_start + ebr.part_size;
                                //actualizar particiones
                                fseek(archivo, TemporalDesplz, SEEK_SET);
                                fwrite(&ebr, sizeof(EBR), 1, archivo);

                                fseek(archivo, desplazamiento, SEEK_SET);
                                fwrite(&ebr, sizeof(EBR), 1, archivo);
                                cout << "Se actualizo la particion con exito" << endl;
                                fclose(archivo);
                                return;
                            } else {
                                cout << "Error: No hay espacio suficiente para aumentar la particion" << endl;
                                fclose(archivo);
                                return;
                            }
                            break;
                        }
                        TemporalDesplz += ebr.part_size + sizeof(EBR)+1;
                    }
                } while (ebr.part_size != 0);
            }
        }
        else{
            string Eliminar = Delete;
            if(Eliminar != "full"){
                cout << "Error: No se reconoce el parametro delete" << endl;
                fclose(archivo);
                return;
            } else{
                string Nombre = Name;
                int indiceParticion = 0;
                if(strcmp(disk.mbr_partition_1.part_name, Nombre.c_str())== 0){
                    indiceParticion = 1;
                } else if(strcmp(disk.mbr_partition_2.part_name, Nombre.c_str())== 0){
                    indiceParticion = 2;
                } else if(strcmp(disk.mbr_partition_3.part_name, Nombre.c_str())== 0){
                    indiceParticion = 3;
                } else if(strcmp(disk.mbr_partition_4.part_name, Nombre.c_str())== 0){
                    indiceParticion = 4;
                }
                if(indiceParticion !=0){
                    Partition particionEliminar;
                    Partition particion;
                    if(indiceParticion == 1){
                        particionEliminar = disk.mbr_partition_1;
                        disk.mbr_partition_1.part_status = '0';
                        disk.mbr_partition_1.part_type = '0';
                        disk.mbr_partition_1.part_fit = '0';
                        disk.mbr_partition_1.part_start = 0;
                        disk.mbr_partition_1.part_size = 0;
                        strcpy(disk.mbr_partition_1.part_name, "");
                        fseek(archivo, particionEliminar.part_start, SEEK_SET);
                        char vacio = '\0';
                        for(int i = 0; i < particionEliminar.part_size; i++){
                            fwrite(&vacio, sizeof(char), 1, archivo);
                        }
                        //actualizar particiones en el mbr
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);
                        cout << "Se elimino la particion con exito" << endl;
                    } else if (indiceParticion ==2){
                        particionEliminar = disk.mbr_partition_2;
                        disk.mbr_partition_2.part_status = '0';
                        disk.mbr_partition_2.part_type = '0';
                        disk.mbr_partition_2.part_fit = '0';
                        disk.mbr_partition_2.part_start = 0;
                        disk.mbr_partition_2.part_size = 0;
                        strcpy(disk.mbr_partition_2.part_name, "");
                        fseek(archivo, particionEliminar.part_start, SEEK_SET);
                        char vacio = '\0';
                        for(int i = 0; i < particionEliminar.part_size; i++){
                            fwrite(&vacio, sizeof(char), 1, archivo);
                        }
                        //actualizar particiones en el mbr
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);
                        cout << "Se elimino la particion con exito" << endl;
                    } else if (indiceParticion ==3){
                        particionEliminar = disk.mbr_partition_3;
                        disk.mbr_partition_3.part_status = '0';
                        disk.mbr_partition_3.part_type = '0';
                        disk.mbr_partition_3.part_fit = '0';
                        disk.mbr_partition_3.part_start = 0;
                        disk.mbr_partition_3.part_size = 0;
                        strcpy(disk.mbr_partition_3.part_name, "");
                        fseek(archivo, particionEliminar.part_start, SEEK_SET);
                        char vacio = '\0';
                        for(int i = 0; i < particionEliminar.part_size; i++){
                            fwrite(&vacio, sizeof(char), 1, archivo);
                        }
                        //actualizar particiones en el mbr
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);
                        cout << "Se elimino la particion con exito" << endl;
                    } else if (indiceParticion ==4){
                        particionEliminar = disk.mbr_partition_4;
                        disk.mbr_partition_4.part_status = '0';
                        disk.mbr_partition_4.part_type = '0';
                        disk.mbr_partition_4.part_fit = '0';
                        disk.mbr_partition_4.part_start = 0;
                        disk.mbr_partition_4.part_size = 0;
                        strcpy(disk.mbr_partition_4.part_name, "");
                        fseek(archivo, particionEliminar.part_start, SEEK_SET);
                        char vacio = '\0';
                        for(int i = 0; i < particionEliminar.part_size; i++){
                            fwrite(&vacio, sizeof(char), 1, archivo);
                        }
                        cout << "Se elimino la particion con exito" << endl;
                        //actualizar particiones en el mbr
                        fseek(archivo, 0, SEEK_SET);
                        fwrite(&disk, sizeof(MBR), 1, archivo);
                    }
                    fclose(archivo);
                } else{
                    Partition particionExtendida;
                    int indiceParticion = 0;

                    if(disk.mbr_partition_1.part_type == 'e'){
                        particionExtendida = disk.mbr_partition_1;
                        indiceParticion = 1;
                    } else if(disk.mbr_partition_2.part_type == 'e'){
                        particionExtendida = disk.mbr_partition_2;
                        indiceParticion = 2;
                    } else if(disk.mbr_partition_3.part_type == 'e'){
                        particionExtendida = disk.mbr_partition_3;
                        indiceParticion = 3;
                    } else if(disk.mbr_partition_4.part_type == 'e'){
                        particionExtendida = disk.mbr_partition_4;
                        indiceParticion = 4;
                    }
                    if(indiceParticion == 0){
                        cout << "Error: No se encontro una particion extendida" << endl;
                        fclose(archivo);
                        return;
                    }
                    EBR ebr;
                    int TemporalDesplz = particionExtendida.part_start ;
                    bool encontrado = false;
                    do {
                        fseek(archivo, TemporalDesplz, SEEK_SET);
                        fread(&ebr, sizeof(EBR), 1, archivo);
                        /*
                         * if(strcmp(ebr.part_name, Nombre.c_str())== 0){
                            encontrado = true;
                            break;
                        }*/
                        if (ebr.part_size != 0){
                            if(strcmp(ebr.part_name, Nombre.c_str())== 0) {
                                encontrado = true;
                                break;
                            }
                            TemporalDesplz += ebr.part_size + sizeof(EBR) + 1;
                        }
                    } while (ebr.part_size != 0);
                    if(encontrado){
                        int dezplAux = TemporalDesplz;
                        TemporalDesplz += sizeof(EBR)+1+ebr.part_size;
                        fseek(archivo, TemporalDesplz, SEEK_SET);
                        fread(&ebr, sizeof(EBR), 1, archivo);

                        do {
                            fseek(archivo, dezplAux, SEEK_SET);
                            fwrite(&ebr, sizeof(EBR), 1, archivo);
                            dezplAux = TemporalDesplz;
                            TemporalDesplz += sizeof(EBR)+1+ebr.part_size;
                            fseek(archivo, TemporalDesplz, SEEK_SET);
                            fread(&ebr, sizeof(EBR), 1, archivo);
                        } while (ebr.part_size != 0);

                        char vacio = '\0';
                        fseek(archivo, dezplAux, SEEK_SET);
                        for (int i = 0; i < sizeof(EBR); i++) {
                            fwrite(&vacio, sizeof(char), 1, archivo);
                        }
                        cout << "Se elimino la particion con exito" << endl;
                        fclose(archivo);
                        return;
                    } else{
                        cout << "Error: No se encontro la particion" << endl;
                        fclose(archivo);
                        return;
                    }
                }
            }
        }
    }
}
